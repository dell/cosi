---
title: ObjectScale
linktitle: ObjectScale
weight: 1
Description: Code features for ObjectScale COSI Driver
---

## Bucket Creation Feature

In order to use COSI Driver on ObjectScale platform, ensure the following omponents are deployed to your cluster:
- Kubernetes Container Object Storage Interface CRDs
- Container Object Storage Interface Controller

### Bucket Class

Installation of ObjectScale COSI driver does not create BucketClass resource. The default sample is shown below:

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

The following is a sample manifest for creating a BucketClaim resource:

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