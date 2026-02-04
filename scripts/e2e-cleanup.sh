#!/usr/bin/env bash
# Copyright © 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
#
# This software contains the intellectual property of Dell Inc.
# or is licensed to Dell Inc. from third parties. Use of this software
# and the intellectual property contained therein is expressly limited to the
# terms and conditions of the License Agreement under which it is provided by or
# on behalf of Dell Inc. or its subsidiaries.

if [ -n "${DEBUG}" ]; then
  set -x
fi

echo "Cleaning up COSI resources"

NS=("access-namespace" "access-grant-namespace" "access-grant-namespace-greenfield" "access-grant-namespace-brownfield" "access-revoke-namespace" "creation-namespace" "deletion-namespace" "access-grant-multiple-namespace-greenfield")

# Cleanup cluster level resources.
for n in "${NS[@]}";
do
    if kubectl get namespace ${n} > /dev/null 2>&1; then
      echo "Cleaning namespace $n"
    else
      echo "Namespace $n does not exist, skipping..."
      continue
    fi

    for s in $(kubectl get secret -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch secret -n=${n} ${s} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    for b in $(kubectl get bucketclaim.objectstorage.k8s.io -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch bucketclaim.objectstorage.k8s.io -n=${n} ${b} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    for b in $(kubectl get bucketaccess.objectstorage.k8s.io -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch bucketaccess.objectstorage.k8s.io -n=${n} ${b} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    for b in $(kubectl get bucket.objectstorage.k8s.io -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch bucket.objectstorage.k8s.io -n=${n} ${b} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    for b in $(kubectl get bucketaccessclass.objectstorage.k8s.io -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch bucketaccessclass.objectstorage.k8s.io -n=${n} ${b} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    for b in $(kubectl get bucketclass.objectstorage.k8s.io -n=${n} -o=jsonpath='{.items[*].metadata.name}');
    do
      kubectl patch bucketclass.objectstorage.k8s.io -n=${n} ${b} -p='{"metadata":{"finalizers":null}}' --type=merge
    done

    # delete all objects from those namespaces
    kubectl delete bucketclaims.objectstorage.k8s.io -n=${n} --all
    kubectl delete bucketaccesses.objectstorage.k8s.io -n=${n} --all
    kubectl delete bucketaccessclasses.objectstorage.k8s.io --all
    kubectl delete bucketclasses.objectstorage.k8s.io --all
    kubectl delete buckets.objectstorage.k8s.io --all
    kubectl delete secret -n=${n} --all
    kubectl delete namespace ${n} --wait --timeout=1m
done

# Cleanup ObjectScale resources in our test namespace.
if [ -n "${OBJECTSCALE_NAMESPACE}" ] && [ -n "${OBJECTSCALE_USER}" ] && [ -n "${OBJECTSCALE_PASSWORD}" ] && [ -n "${OBJECTSCALE_GATEWAY}" ]; then
    echo "Cleaning up ObjectScale resources in ${OBJECTSCALE_NAMESPACE}"

    AUTH_TOKEN=$(curl -ski -u "${OBJECTSCALE_USER}:${OBJECTSCALE_PASSWORD}" -X GET "${OBJECTSCALE_GATEWAY}/login" | grep X-SDS-AUTH-TOKEN)
    if [ -z "$AUTH_TOKEN" ]; then
        echo "ERROR: Failed to get authorization token"
        exit 1
    fi

    buckets=$(curl -sk -H "${AUTH_TOKEN}" -H "Accept: application/json" -X GET "${OBJECTSCALE_GATEWAY}/object/bucket" | jq -r '.object_bucket[] | .name')
    if [ ${PIPESTATUS[0]} -ne 0 ]; then
        echo "ERROR: Failed to get buckets from namespace ${OBJECTSCALE_NAMESPACE}"
        exit 1
    fi

    for bucket in $buckets
    do
        echo "Deleting bucket ${bucket}"
        status_code=$(curl -sk -w "%{http_code}" -H "${AUTH_TOKEN}" -H "Accept: application/json" -X POST "${OBJECTSCALE_GATEWAY}/object/bucket/${bucket}/deactivate?emptyBucket=true")
        if [ "$status_code" != "200" ] && [ "$status_code" != "202" ]; then
            echo "WARNING: Failed to delete bucket ${bucket} https status code ${status_code}"
        fi
    done
fi

echo "Done."
