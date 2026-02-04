#!/usr/bin/env bash
# Copyright © 2023-2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# This software contains the intellectual property of Dell Inc.
# or is licensed to Dell Inc. from third parties. Use of this software
# and the intellectual property contained therein is expressly limited to the
# terms and conditions of the License Agreement under which it is provided by or
# on behalf of Dell Inc. or its subsidiaries.

if [ -n "${DEBUG}" ]; then
  set -x
fi

#########################################################################################
# Configuration:
# - The shell shall write a message to standard error when it tries to expand a variable
#   that is not set and immediately exit.
#----------------------------------------------------------------------------------------
set -u

# Helm specific
export DRIVER_NAMESPACE="${DRIVER_NAMESPACE:-cosi-test-ns}"
export HELM_RELEASE_NAME="${HELM_RELEASE_NAME:-dell-cosi}"

# Image specific
export REGISTRY="${REGISTRY:-quay.io/dell/container-storage-modules}"
export IMAGENAME="${IMAGENAME:-cosi:nightly}"
export CHART_BRANCH="${CHART_BRANCH:-main}"

# ObjectScale specific
export OBJECTSCALE_NAMESPACE="${OBJECTSCALE_NAMESPACE}"
export OBJECTSCALE_USER="${OBJECTSCALE_USER}"
export OBJECTSCALE_PASSWORD="${OBJECTSCALE_PASSWORD}"
export OBJECTSCALE_GATEWAY="${OBJECTSCALE_GATEWAY}"
export OBJECTSCALE_S3_ENDPOINT="${OBJECTSCALE_S3_ENDPOINT}"

# Tests specific
export DRIVER_CONTAINER_NAME="${DRIVER_CONTAINER_NAME:-objectstorage-provisioner}"

#########################################################################################
# Main:
# - subshell execution
#----------------------------------------------------------------------------------------
(
./e2e-cleanup.sh
if [ $? -ne 0 ]; then
  echo "ERROR: e2e-cleanup.sh failed"
  exit 1
fi

# uninstall driver
helm uninstall "${HELM_RELEASE_NAME}" -n="${DRIVER_NAMESPACE}" || true
kubectl delete leases -n="${DRIVER_NAMESPACE}" cosi-dellemc-com-cosi || true

# save driver configuration values in a file
cat > /tmp/cosi-conf.yml <<EOF
connections:
- objectscale:
    id: e2e.test.objectscale
    namespace: ${OBJECTSCALE_NAMESPACE}
    credentials:
      username: ${OBJECTSCALE_USER}
      password: ${OBJECTSCALE_PASSWORD}
    mgmt-endpoint: ${OBJECTSCALE_GATEWAY}
    region: us-east-1
    emptyBucket: false
    protocols:
      s3:
        endpoint: ${OBJECTSCALE_S3_ENDPOINT}
    tls:
      insecure: true
EOF

cosiConfig=`cat /tmp/cosi-conf.yml | base64 -w 0`

kubectl apply -f - <<EOF
apiVersion: v1
data:
  config.yaml: ${cosiConfig}
kind: Secret
metadata:
  name: ${HELM_RELEASE_NAME}-config
  namespace: ${DRIVER_NAMESPACE}
type: Opaque
EOF

# When this option is on, if a simple command fails for any of the reasons listed in
# Consequences of Shell Errors or returns an exit status value >0, and is not part of the
# compound list following a while, until, or if keyword, and is not a part of an AND or
# OR list, and is not a pipeline preceded by the ! reserved word, then the shell shall
# immediately exit.
set -e

rm -rf helm
git clone \
  --branch "${CHART_BRANCH}" \
  --single-branch \
  https://github.com/dell/helm-charts.git helm

# install the driver
helm install "${HELM_RELEASE_NAME}" ./helm/charts/cosi \
  --set=images.provisioner.image="${REGISTRY}/${IMAGENAME}" \
  --set=imagePullPolicy=Always \
  --namespace="${DRIVER_NAMESPACE}" \
  --create-namespace

# check if the driver is installed correctly
kubectl wait \
  --for=condition=available \
  --timeout=60s \
  --namespace="${DRIVER_NAMESPACE}" \
  deployments "${HELM_RELEASE_NAME}"

# start e2e tests
cd .. && make integration-test
)
