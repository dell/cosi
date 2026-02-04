#!/bin/bash
#
#Copyright © 2020-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
PROG="${0}"
source "$SCRIPTDIR"/common.sh
DRIVER="cosi"

if [ -z "${DEBUGLOG}" ]; then
  export DEBUGLOG="${SCRIPTDIR}/install-debug.log"
  if [ -f "${DEBUGLOG}" ]; then
    rm -f "${DEBUGLOG}"
  fi
fi

declare -a VALIDDRIVERS


# source-verify-driver will call the proper method to source verification method scripts
function source-verify-driver() {
  if [ -z "${1}" ]; then
    decho "Expected one argument: 'the driver name' to verify-driver. Received none."
    exit $EXIT_ERROR
  fi
  local D="${1}"

  # check if a script, matching a specific name, exists
  local SCRIPTNAME="verify-$D.sh"
  # check if the script exists, and source it
  if [ -f "${SCRIPTDIR}/${SCRIPTNAME}" ]; then
    source "${SCRIPTDIR}/${SCRIPTNAME}"
  else
    echo "Unable to find ${SCRIPTDIR}/${SCRIPTNAME}"
  fi
}


#
# verify-driver will call the proper method to verify a specific driver
function verify-driver() {
  if [ -z "${1}" ]; then
    decho "Expected one argument: 'the driver name' to verify-driver. Received none."
    exit $EXIT_ERROR
  fi
  local D="${1}"
  # source a script containing verification methods
  source-verify-driver "${D}"
  # check if a verify-$DRIVER function exists
  # if not, error and exit
  # if yes, check to see if it should be run and run it
  FNTYPE=$(type -t verify-$D)
  if [ "$FNTYPE" != "function" ]; then
    decho "ERROR: verify-$D function does not exist"
    exit $EXIT_ERROR
  else
    header
    log step "Driver: ${D}"
    decho
    verify-$D
    summary
  fi
}

# Print usage information
function usage() {
  decho
  decho "Help for $PROG"
  decho
  decho "Usage: $PROG options..."
  decho "Options:"
  decho "  Required"
  decho "  --namespace[=]<namespace>       Kubernetes namespace to install the COSI driver"
  decho "  --values[=]<values.yaml>        Values file, which defines configuration values"
  decho "  --driver-version                verify COSI Driver Image version is in values file"

  decho "  Optional"
  decho "  --release[=]<helm release>      Name to register with helm, default value will match the driver name"
  decho "  -h                              Help"
  decho

  exit $EXIT_WARNING
}

# print header information
function header() {
  log section "Verifying Kubernetes and driver configuration"
  echo "|- Kubernetes Version: ${kMajorVersion}.${kMinorVersion}"
}

# verify secrets exist
function verify_required_secrets() {
  log step "Verifying that required secrets have been created"

  error=0
  for N in "${@}"; do
    # Make sure the secret has already been established
    run_command kubectl get secrets -n "${NS}" 2>/dev/null | grep "${N}" --quiet
    if [ $? -ne 0 ]; then
      error=1
      found_error "Required secret, ${N}, does not exist."
    fi
  done
  check_error error
}

# verify minimum and maximum k8s versions
function verify_k8s_versions() {
  if [ "${OPENSHIFT}" == "true" ]; then
    return
  fi
  log step "Verifying Kubernetes version"
  decho

  local MIN=${1}
  local MAX=${2}
  local V="${kMajorVersion}.${kMinorVersion}"
  # check non supported version (k8s alpha/beta)
  if [ -n "${kNonGAVersion}" ]; then
    echo "Installing on an unreleased version of Kubernetes : "${kNonGAVersion}". Acknowledge and proceed with installation? (y/n)"
    read -n 1 -p "Press 'y' to continue or any other key to exit: " CONT
    decho
      if [ "${CONT}" != "Y" -a "${CONT}" != "y" ]; then
        decho "quitting at user request"
        exit 2
      fi
   fi

  # check minimum
  log arrow
  log smart_step "Verifying minimum Kubernetes version" "small"
  error=0
  if [[ ${V} < ${MIN} ]]; then
    error=1
    found_warning "Kubernetes version ${V} is too old. Minimum required version is: ${MIN}"
  fi
  check_error error

  # check maximum
  log arrow
  log smart_step "Verifying maximum Kubernetes version" "small"
  error=0
  if [[ ${V} > ${MAX} ]]; then
    error=1
    found_warning "Kubernetes version ${V} is newer than the version that has been tested. Latest tested version is: ${MAX}"
  fi
  check_error error

}

# verify minimum and maximum openshift versions
function verify_openshift_versions() {
  if [ "${OPENSHIFT}" != "true" ]; then
    return
  fi
  log step "Verifying OpenShift version"
  decho

  local MIN=${1}
  local MAX=${2}
  local V=$(OpenShiftVersion)
  # check minimum
  log arrow
  log smart_step "Verifying minimum OpenShift version" "small"
  error=0
  if (( ${V%%.*} < ${MIN%%.*} || ( ${V%%.*} == ${MIN%%.*} && ${V##*.} < ${MIN##*.} ) )) ; then
    error=1
    found_warning "OpenShift version ${V} is too old. Minimum required version is: ${MIN}"
  fi
  check_error error

  # check maximum
  log arrow
  log smart_step "Verifying maximum OpenShift version" "small"
  error=0
  if (( ${V%%.*} > ${MAX%%.*} || ( ${V%%.*} == ${MAX%%.*} && ${V##*.} > ${MAX##*.} ) )) ; then
    error=1
    found_warning "OpenShift version ${V} is newer than the version that has been tested. Latest tested version is: ${MAX}"
  fi
  check_error error
}

# verify namespace
function verify_namespace() {
  log step "Verifying that required namespaces have been created"

  error=0
  for N in "${@}"; do
    # Make sure the namespace exists
    run_command kubectl describe namespace "${N}" >/dev/null 2>&1
    if [ $? -ne 0 ]; then
      error=1
      found_error "Namespace does not exist: ${N}"
    fi
  done

  check_error error
}

# verify that helm is v3 or above
function verify_helm_3() {
  log step "Verifying helm version"

  error=0
  # Check helm installer version
  helm --help >&/dev/null || {
    found_error "helm is required for installation"
    log step_failure
    return
  }

  run_command helm version | grep "v3." --quiet
  if [ $? -ne 0 ]; then
    error=1
    found_error "Driver installation is supported only using helm 3"
  fi

  check_error error
}


# found_error, installation will not continue
function found_error() {
  for N in "$@"; do
    ERRORS+=("${N}")
  done
}

# found_warning, installation can continue
function found_warning() {
  for N in "$@"; do
    WARNINGS+=("${N}")
  done
}

# Print a nice summary at the end
function summary() {
  local VERSTATUS="Success"
  if [ "${#WARNINGS[@]}" -ne 0 ]; then
    VERSTATUS="With Warnings"
  fi
  if [ "${#ERRORS[@]}" -ne 0 ]; then
    VERSTATUS="With Errors"
  fi
  decho
  log section "Verification Complete - ${VERSTATUS}"
  # print all the WARNINGS
  if [ "${#WARNINGS[@]}" -ne 0 ]; then
    log warnings
    for E in "${WARNINGS[@]}"; do
      decho "- ${E}"
    done
    RC=$EXIT_WARNING
  fi

  # print all the ERRORS
  if [ "${#ERRORS[@]}" -ne 0 ]; then
    log errors
    for E in "${ERRORS[@]}"; do
      decho "- ${E}"
    done
    RC=$EXIT_ERROR
  fi

  return $RC
}

#
# validate_params will validate the parameters passed in
function validate_params() {
  # make sure the driver was specified
  if [ -z "${DRIVER}" ]; then
    decho "No driver specified"
    usage
    exit 1
  fi
  # the namespace is required
  if [ -z "${NS}" ]; then
    decho "No namespace specified"
    usage
    exit 1
  fi
  # values file
  if [ -z "${VALUES}" ]; then
    decho "No values file was specified"
    usage
    exit 1
  fi
  if [ ! -f "${VALUES}" ]; then
    decho "Unable to read values file at: ${VALUES}"
    usage
    exit 1
  fi
}

# verify helm values file version matches helm chart version, which is the driver version
# args: expected values version
function verify_helm_values_version() {
  log step "Verifying helm values version"

  error=0
  values_version=$(grep ^version: $VALUES | awk '{print $2}' | tr -d '"')
  # if values version is not found, it's ok
  if [ -z $values_version ]; then
    error=1
    found_error "Helm values file version not found"
    check_error error
    return
  fi

  match=false
  for expected_version in "$@"; do
    if [ "$expected_version" = "$values_version" ]; then
      match=true
      break
    fi
  done

  if [ "$match" = false ]; then
    error=1
    found_error "Incompatible helm values file specified - expected: `echo $@| tr " " "/ "`, found: $values_version"
  fi

  check_error error
}

#
# main
#
# default values

# exit codes
EXIT_SUCCESS=0
EXIT_WARNING=1
EXIT_ERROR=99

# arrays of messages
WARNINGS=()
ERRORS=()

# make sure kubectl is available
kubectl --help >&/dev/null || {
  decho "kubectl required for verification... exiting"
  exit $EXIT_ERROR
}

# Determine the nodes
MINION_NODES=$(run_command kubectl get nodes -o wide | grep -v -e master -e INTERNAL | awk ' { print $6; }')
MASTER_NODES=$(run_command kubectl get nodes -o wide | awk ' /master/{ print $6; }')
# Get the kubernetes major and minor version numbers.
kMajorVersion=$(run_command kubectl version | grep 'Server Version' | sed -E 's/.*v([0-9]+)\.[0-9]+\.[0-9]+.*/\1/')
kMinorVersion=$(run_command kubectl version | grep 'Server Version' | sed -E 's/.*v[0-9]+\.([0-9]+)\.[0-9]+.*/\1/')
kNonGAVersion=$(run_command kubectl version | grep 'Server Version' | sed -n 's/.*\(-[alpha|beta][^ ]*\).*/\1/p')

while getopts ":h-:" optchar; do
  case "${optchar}" in
  -)
    case "${OPTARG}" in
      # NAMESPACE
    namespace)
      NS="${!OPTIND}"
      if [[ -z ${NS} || ${NS} == "--skip-verify" ]]; then
        NS=${DEFAULT_NS}
      else
        OPTIND=$((OPTIND + 1))
      fi
      ;;
    namespace=*)
      NS=${OPTARG#*=}
      if [[ -z ${NS} ]]; then NS=${DEFAULT_NS}; fi
      ;;
      # RELEASE
    version)
      VERSION="${!OPTIND}"
      OPTIND=$((OPTIND + 1))
      ;;
    driver-version)
      DRIVER_VERSION="${!OPTIND}"
      OPTIND=$((OPTIND + 1))
      ;;
    release)
      RELEASE="${!OPTIND}"
      OPTIND=$((OPTIND + 1))
      ;;
    release=*)
      RELEASE=${OPTARG#*=}
      ;;
      # VALUES
    values)
      VALUES="${!OPTIND}"
      OPTIND=$((OPTIND + 1))
      ;;
    values=*)
      VALUES=${OPTARG#*=}
      ;;
    *)
      decho "Unknown option --${OPTARG}"
      decho "For help, run $PROG -h"
      exit $EXIT_ERROR
      ;;
    esac
    ;;
  h)
    usage
    ;;
  *)
    decho "Unknown option -${OPTARG}"
    decho "For help, run $PROG -h"
    exit $EXIT_ERROR
    ;;
  esac
done

# by default the NAME of the helm release of the driver is the same as the driver name
RELEASE=$(get_release_name "${DRIVER}")

# validate the parameters passed in
validate_params "${MODE}"
OPENSHIFT=$(isOpenShift)

verify-driver "${DRIVER}"
exit $?
