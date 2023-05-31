@component_COSI

Feature: Bucket creation on ObjectScale platform

    As an ObjectScale platform user
    I want to add BucketClaim, which is a request for a new Bucket
    so that the information about new Bucket (e.g. BucketID) are returned

    Background:
        Given Kubernetes cluster is up and running
        And ObjectScale platform is installed on the cluster
        And ObjectStore "${objectstoreName}" is created
        And Kubernetes namespace "cosi-driver" is created
        And Kubernetes namespace "namespace-1" is created
        And COSI controller "objectstorage-controller" is installed in namespace "default"
        And COSI driver "cosi-driver" is installed in namespace "cosi-driver"
        And specification of custom resource "my-bucket-class" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketClass
        metadata:
            name: my-bucket-class
        deletionPolicy: delete
        driverName: cosi.dellemc.com
        parameters:
            ID: ${driverID}
        """
        And specification of custom resource "bucket-claim-valid" is:
        """
        apiVersion: v1
        kind: BucketClaim
        metadata:
            name: bucket-claim-valid
            namespace: namespace-1
        spec:
            bucketClassName: my-bucket-class
            protocol: S3
        """
        And specification of custom resource "bucket-claim-invalid" is:
        """
        apiVersion: v1
        kind: BucketClaim
        metadata:
            name: bucket-claim-invalid
            namespace: namespace-1
        spec:
            bucketClassName: bucket-class-invalid
            protocol: S3
        """
        And BucketClass resource is created from specification "my-bucket-class"
    
    Scenario: Successfull bucket creation
        When BucketClaim resource is created from specification "bucket-claim-valid"
        And Bucket resource referencing BucketClaim resource "bucket-claim-valid" is created
        Then Bucket resource referencing BucketClaim resource "bucket-claim-valid" is created in ObjectStore "${objectstoreName}"
        And BucketClaim resource "bucket-claim-valid" in namespace "namespace-1" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-valid" bucketID is not empty

    Scenario: Unsuccessfull bucket creation
        When BucketClaim resource is created from specification "bucket-claim-invalid"
        And Bucket resource referencing BucketClaim resource "bucket-claim-invalid" is created
        Then Bucket resource referencing BucketClaim resource "bucket-claim-invalid" is not created in ObjectStore "${objectstoreName}"
        And BucketClaim resource "bucket-claim-invalid" in namespace "namespace-1" status "bucketReady" is "false"
        And BucketClaim events contains an error: "Cannot create Bucket: BucketClass does not exist"
