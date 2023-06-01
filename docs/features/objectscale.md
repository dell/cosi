---
title: ObjectScale
linktitle: ObjectScale
weight: 1
Description: Code features for ObjectScale COSI Driver
---

<!--
TODO: add brownfield provisioning info to Bucket
TODO: add brownfield provisioning info to BucketClaim
TODO: add description of `spec.existingBucketName` for Bucket
TODO: add description of `spec.existingBucketName` for BucketClaim
TODO: add description of `spec.authenticationType` for BucketAccessClass
TODO: write Bucket Access Granting feature description
FIXME: is the `parameters.driverID` a good name?
-->

In order to use COSI Driver on ObjectScale platform, ensure the following components are deployed to your cluster:
- Kubernetes Container Object Storage Interface CRDs
- Container Object Storage Interface Controller

## Notational Conventions

The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119](http://tools.ietf.org/html/rfc2119) (Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels", BCP 14, RFC 2119, March 1997).

Fields are specified by their's path. Consider the following examples:

1. Field specified by the following path `spec.authenticationType=IAM` is reflected in their resouces YAML as the following:

```yaml
spec:
  authenticationType: IAM
```
2. field specified by path `spec.protocols=[Azure,GCS]` is reflected in their resouces YAML as the following:

```yaml
spec:
  protocols:
    - Azure
    - GCS
```

## Bucket Creation Feature

### Bucket

`Bucket` represents a Bucket or its equivalent in the storage backend. Generally, it should be created only in the brownfield provisioning scenario. The following is a sample manifest of `Bucket` resource:

```yaml
apiVersion: objectstorage.k8s.io/v1alpha1
kind: Bucket
metadata:
  name: my-bucket
spec:
  driverName: cosi.dellemc.com
  bucketClassName: my-bucket-class
  bucketClaim: my-bucket-claim
  deletionPolicy: Delete
  protocols:
    - S3
```

### Bucket Claim

`BucketClaim` represents a claim to provision a `Bucket`. The following is a sample manifest for creating a BucketClaim resource:

```yaml
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketClaim
metadata:
  name: my-bucket-claim
  namespace: my-namespace
spec:
  bucketClassName: my-bucket-class
  protocols:
    - S3
```

### Unsupported options

- `spec.protocols=[Azure,GCS]` - Protocols are the set of data API this bucket is required to support. From protocols specified by COSI (`v1alpha1`), Dell ObjectScale platform only supports the S3 protocol. Protocols `Azure` and `GCS` MUST NOT be used.

### Bucket Class

Installation of ObjectScale COSI driver does not create `BucketClass` resource. `BucketClass` represents a class of `Buckets` with similar characteristics. 
Dell COSI Driver is a multi-backend driver, meaning that for every platform the specific `BucketClass` should be created. The `BucketClass` resource should contain the name of multi-backend driver and driverID for specific Object Storage Platform. 
The default sample is shown below:

```yaml
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketClass
metadata:
  name: my-bucket-class
driverName: cosi.dellemc.com
deletionPolicy: Delete
parameters:
  driverID: "objectscale.secure.panda"
```

## Bucket Deletion Feature

There are a few crucial details regarding bucket deletion. The first one is Deletion Policy which is used to specify how COSI should handle deletion of a bucket. It is found in K8s CRD and can be set to Delete and Retain. The second crucial detail is `emptyBucket` field in the Helm Chart configuration.

### `deletionPolicy`

> ⚠ **WARNING**: this field is case sensitive, and the bucket deletion will fail if policy is not set exactly to *Delete* or *Retain*.

DeletionPolicy in `BucketClass` resource is used to specify how COSI should handle deletion of the bucket. There are two possible values:
- **Retain**: Indicates that the bucket should not be deleted from the Object Storage Platform (OSP), it means that the underlying bucket is not cleaned up when the `Bucket` object is deleted. It makes the bucket unreachable from k8s level. 
- **Delete**: Indicates that the bucket should be permanently deleted from the Object Storage Platform (OSP) once all the workloads accessing this bucket are done, it means that the underlying bucket is cleaned up when the Bucket object is deleted.

### `emptyBucket`

`emptyBucket` field is set in config `.yaml` file passed to the chart during COSI driver installation. If it is set to `true`, then the bucket will be emptied before deletion. If it is set to `false`, then Objectscale will not be able to delete not empty bucket and return error.

`emptyBucket` has no effect when Deletion Policy is set to `Retain`.

## Bucket Access Granting Feature

### Bucket Access Class

```yaml
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketAccessClass
metadata:
  name: my-bucket-access-class
driverName: cosi.dellemc.com
authenticationType: KEY
parameters:
  driverID: "objectscale.secure.panda"
```

### `authenticationType`

### Unsupported options

- `authenticationType=IAM` - denotes the style of authentication. The `IAM` value MUST NOT be used, because IAM style authentication is not supported. 


### Bucket Access

```yaml
apiVersion: objectstorage.k8s.io/v1alpha1
kind: BucketAccess
metadata:
  name: my-bucket-access
  namespace: my-namespace
spec:
  bucketClaimName: my-bucket-claim
  protocol: S3
  bucketAccessClassName: my-bucket-access-class
  credentialsSecretName: my-cosi-secret
```

### `spec.protocol`

> ⚠ **WARNING**: this field is case sensitive, and the provisioning will fail if protocol is not set exactly to *S3*.


### Unsupported options

- `spec.serviceAccountName=...` - is the name of the serviceAccount that COSI will map to the object storage provider service account when IAM styled authentication is specified. As the IAM style authentication is not supported, this field is also unsupported.
- `spec.protocol=...` - Protocols are the set of data API this bucket is required to support. From protocols specified by COSI (`v1alpha1`), Dell ObjectScale platform only supports the S3 protocol. Protocols `Azure` and `GCS` MUST NOT be used.

## Bucket Access Revoking Feature
