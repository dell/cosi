---
title: ObjectScale
linktitle: ObjectScale
weight: 1
Description: Code features for ObjectScale COSI Driver
---

## Bucket Creation Feature

In order to use COSI Driver on ObjectScale platform, ensure the following components are deployed to your cluster:
- Kubernetes Container Object Storage Interface CRDs
- Container Object Storage Interface Controller

### Bucket Class

Installation of ObjectScale COSI driver does not create `BucketClass` resource. `BucketClass` represents a class of `Buckets` with similar characteristics. 
Dell COSI Driver is a multi-backend driver, meaning that for every platform the specific `BucketClass` should be created. The `BucketClass` resource should contain the name of multi-backend driver and driverID for specific Object Storage Platform. 
The default sample is shown below:

```yaml
kind: BucketClass
apiVersion: objectstorage.k8s.io/v1alpha1
metadata:
  name: my-bucket-class
driverName: dell.objectstorage.k8s.io
deletionPolicy: Delete
parameters:
  driverID: "objectscale.secure.panda"
```

### Bucket Claim

`BucketClaim` represents a claim to provision a `Bucket`. The following is a sample manifest for creating a BucketClaim resource:

```yaml
kind: BucketClaim
apiVersion: objectstorage.k8s.io/v1alpha1
metadata:
  name: my-bucket-claim
spec:
  bucketClassName: my-bucket-class
  protocols:
    - S3
```

### Bucket

`Bucket` represents a Bucket or its equivalent in the storage backend. The following is a sample manifest of `Bucket` resource:

```yaml
kind: Bucket
apiVersion: objectstorage.k8s.io/v1alpha1
metadata:
  name: my-bucket
spec:
  bucketClassName: my-bucket-class
  bucketClaimName: my-bucket-claim
  driverName: dell.objectstorage.k8s.io
  deletePolicy: delete
  protocols:
    - S3
```

## Bucket Deletion Feature

There are a few crucial details regarding bucket deletion. The first one is Deletion Policy which is used to specify how COSI should handle deletion of a bucket. It is found in K8s CRD and can be set to Delete and Retain. The second crucial detail is `emptyBucket` field in the Helm Chart configuration.

### Deletion Policy

DeletionPolicy in `BucketClass` resource is used to specify how COSI should handle deletion of the bucket. There are two possible values: 
- **Retain**: Indicates that the bucket should not be deleted from the Object Storage Platform (OSP), it means that the underlying bucket is not cleaned up when the `Bucket` object is deleted. It makes the bucket unreachable from k8s level. 
- **Delete**: Indicates that the bucket should be permanently deleted from the Object Storage Platform (OSP) once all the workloads accessing this bucket are done, it means that the underlying bucket is cleaned up when the Bucket object is deleted.

### emptyBucket 

`emptyBucket` field is set in config `.yaml` file passed to the chart during COSi driver installation. If it is set to `true`, then the bucket will be emptied before deletion. If it is set to `false`, then Objectscale will not be able to delete not empty bucket and return error.

`emptyBucket` has no effect when Deletion Policy is set to `Retain`.

## Bucket Access Granting Feature

## Bucket Access Revoking Feature