#!/usr/bin/env bash

# subshell execution
(

NS=("access-namespace" "access-grant-namespace" "access-revoke-namespace" "creation-namespace" "deletion-namespace")

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
  for s in $(kubectl get secret -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch secret -n "${n}" "${s}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketclaim.objectstorage.k8s.io -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketclaim.objectstorage.k8s.io -n "${n}" "{$b}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketaccess.objectstorage.k8s.io -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketaccess.objectstorage.k8s.io -n "${n}" "{$b}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucket.objectstorage.k8s.io -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucket.objectstorage.k8s.io -n "${n}" "{$b}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketaccessclass.objectstorage.k8s.io -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketaccessclass.objectstorage.k8s.io -n "${n}" "{$b}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  for b in $(kubectl get bucketclass.objectstorage.k8s.io -n "${n}" -o jsonpath='{.items[*].metadata.name}');
  do
    kubectl patch bucketclass.objectstorage.k8s.io -n "${n}" "${b}" -p '{"metadata":{"finalizers":null}}' --type=merge
  done

  # delete all objects from those namespaces
  kubectl delete bucketclaims.objectstorage.k8s.io -n="${n}" --all
  kubectl delete bucketaccesses.objectstorage.k8s.io -n="${n}" --all
  kubectl delete bucketaccessclasses.objectstorage.k8s.io --all
  kubectl delete bucketclasses.objectstorage.k8s.io --all
  kubectl delete buckets.objectstorage.k8s.io --all
  kubectl delete secret -n "${n}" --all
  kubectl delete namespace "${n}" 
done

# uninstall driver
helm uninstall cosi-driver -n="$DRIVER_NAMESPACE" || true
kubectl delete leases -n="$DRIVER_NAMESPACE" cosi-dellemc-com-cosi || true

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

# go to the cosi-driver folder 
cd "${DRIVER_REPO_PATH}" || exit 1

# install the driver
helm install cosi-driver ./helm/cosi-driver \
--set provisioner.image.repository="${REGISTRY}"/cosi-driver \
--set provisioner.image.tag="$(git rev-parse HEAD)" \
--set provisioner.image.pullPolicy=Always \
--set provisioner.logLevel=trace \
--set sidecar.verbosity=low \
--set=provisioner.logFormat=json \
--set-file configuration.data=/tmp/cosi-conf.yml \
--namespace=cosi-driver \
--create-namespace

# check if the driver is installed correctly 
kubectl wait --for=condition=available --timeout=60s deployment/cosi-driver -n="${DRIVER_NAMESPACE}"

# start e2e tests
make integration-test
)
