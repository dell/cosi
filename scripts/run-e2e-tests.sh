#!/usr/bin/env bash

if [ -n "${DEBUG}" ]; then
  # The  shell shall write to standard error a trace for each command after it expands
  # the command and before it executes it.
  set -x
fi

#########################################################################################
# Configuration:
# - The shell shall write a message to standard error when it tries to expand a variable
#   that is  not set and immediately exit.
#----------------------------------------------------------------------------------------
set -u

# Helm specific
export DRIVER_NAMESPACE="${DRIVER_NAMESPACE:-cosi-test-ns}"
export HELM_RELEASE_NAME="${HELM_RELEASE_NAME:=dell-cosi}"

# Image specific
export REGISTRY="${REGISTRY:-docker.io}"
export IMAGENAME="${IMAGENAME:-dell/cosi}"
export CHART_BRANCH="${CHART_BRANCH:-main}"

# ObjectScale specific
export OBJECTSCALE_NAMESPACE="${OBJECTSCALE_NAMESPACE}"
export OBJECTSCALE_ID="${OBJECTSCALE_ID}"
export OBJECTSCALE_OBJECTSTORE_ID="${OBJECTSCALE_OBJECTSTORE_ID}"
export OBJECTSCALE_USER="${OBJECTSCALE_USER}"
export OBJECTSCALE_PASSWORD="${OBJECTSCALE_PASSWORD}"
export OBJECTSCALE_GATEWAY="${OBJECTSCALE_GATEWAY}"
export OBJECTSCALE_OBJECTSTORE_GATEWAY="${OBJECTSCALE_OBJECTSTORE_GATEWAY}"
export OBJECTSCALE_S3_ENDPOINT="${OBJECTSCALE_S3_ENDPOINT}"

# Tests specific
export DRIVER_CONTAINER_NAME="${DRIVER_CONTAINER_NAME:-objectstorage-provisioner}"

#########################################################################################
# Main:
# - subshell execution
#----------------------------------------------------------------------------------------
(

NS=("access-namespace" "access-grant-namespace" "access-grant-namespace-greenfield" "access-grant-namespace-brownfield" "access-revoke-namespace" "creation-namespace" "deletion-namespace")

# delete all finalizers and then objects from those namespaces
for n in "${NS[@]}";
do
  # first check if namespace exists
  if kubectl get namespace "${n}" > /dev/null 2>&1; then
    echo "Cleaning namespace $n"
  else
    echo "Namespace $n does not exist, skipping..."
    continue
  fi

  # delete all finalizers and then objects from those namespaces
  for s in $(kubectl get secret -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch secret -n="${n}" "${s}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketclaim.objectstorage.k8s.io -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketclaim.objectstorage.k8s.io -n="${n}" "{$b}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketaccess.objectstorage.k8s.io -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketaccess.objectstorage.k8s.io -n="${n}" "{$b}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucket.objectstorage.k8s.io -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucket.objectstorage.k8s.io -n="${n}" "{$b}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketaccessclass.objectstorage.k8s.io -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketaccessclass.objectstorage.k8s.io -n="${n}" "{$b}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketclass.objectstorage.k8s.io -n="${n}" -o=jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketclass.objectstorage.k8s.io -n="${n}" "${b}" -p='{"metadata":{"finalizers":null}}' --type=merge
  done

  # delete all objects from those namespaces
  kubectl delete bucketclaims.objectstorage.k8s.io -n="${n}" --all
  kubectl delete bucketaccesses.objectstorage.k8s.io -n="${n}" --all
  kubectl delete bucketaccessclasses.objectstorage.k8s.io --all
  kubectl delete bucketclasses.objectstorage.k8s.io --all
  kubectl delete buckets.objectstorage.k8s.io --all
  kubectl delete secret -n="${n}" --all
  kubectl delete namespace "${n}"
done

# uninstall driver
helm uninstall "${HELM_RELEASE_NAME}" -n="${DRIVER_NAMESPACE}" || true
kubectl delete leases -n="${DRIVER_NAMESPACE}" cosi-dellemc-com-cosi || true

# save driver configuration values in a file
cat > /tmp/cosi-conf.yml <<EOF
connections:
- objectscale:
    id: e2e.test.objectscale
    namespace: ${OBJECTSCALE_NAMESPACE}
    objectscale-id: ${OBJECTSCALE_ID}
    objectstore-id: ${OBJECTSCALE_OBJECTSTORE_ID}
    credentials:
      username: ${OBJECTSCALE_USER}
      password: ${OBJECTSCALE_PASSWORD}
    objectscale-gateway: ${OBJECTSCALE_GATEWAY}
    objectstore-gateway: ${OBJECTSCALE_OBJECTSTORE_GATEWAY}
    region: us-east-1
    emptyBucket: false
    protocols:
      s3:
        endpoint: ${OBJECTSCALE_S3_ENDPOINT}
    tls:
      insecure: true
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
  --set=provisioner.image.repository="${REGISTRY}/${IMAGENAME}" \
  --set=provisioner.image.tag="$(git rev-parse HEAD)" \
  --set=provisioner.image.pullPolicy=Always \
  --set=provisioner.logFormat=json \
  --set=provisioner.logLevel=0 \
  --set=provisioner.otelEndpoint='' \
  --set=sidecar.verbosity=10 \
  --set-file=configuration.data=/tmp/cosi-conf.yml \
  --namespace="${DRIVER_NAMESPACE}" \
  --create-namespace

# check if the driver is installed correctly 
kubectl wait \
  --for=condition=available \
  --timeout=60s \
  --namespace="${DRIVER_NAMESPACE}" \
  deployments "${HELM_RELEASE_NAME}"

# start e2e tests
make integration-test
)
