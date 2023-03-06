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

## Bucket Deletion Fetaure

## Bucket Access Granting Feature

## Bucket Access Revoking Feature