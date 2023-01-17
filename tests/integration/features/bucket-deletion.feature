@component_COSI

Feature: Bucket deletion from ObjectScale platform

    As an ObjectScale platform user
    I want to delete BucketClaim
    so that existing Bucket is deleted or left regarding to deletionPolicy

    Background: 
        Given Kubernetes cluster is up and running
        And ObjectScale platform is installed on the cluster
        And ObjectStore "object-store-1" is created
        And Kubernetes namespace "driver-ns" is created
        And Kubernetes namespace "namespace-1" is created
        And COSI controller "cosi-controller" is installed in namespace "driver-ns"
        And COSI driver "cosi-driver" is installed in namespace "driver-ns"

    Scenario: BucketClaim deletion with deletionPolicy set to "delete"
        Given specification of custom resource "my-bucket-class-delete" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketClass
        metadata:
            name: my-bucket-class-delete
        deletionPolicy: delete
        driverName: cosi-driver
        parameters:
            objectScaleID: ${objectScaleID}
            objectStoreID: ${objectStoreID}
            accountSecret: ${secretName}                   
        """
        And specification of custom resource "my-bucket-claim-delete" is:
        """
        apiVersion: v1
        kind: BucketClaim
        metadata:
            name: my-bucket-claim-delete
            namespace: namespace-1
        spec:                                            
            bucketClassName: my-bucket-class-delete
            protocol: S3
        """  
        And BucketClass resource is created from specification "my-bucket-class-delete"
        And BucketClaim resource is created from specification "my-bucket-claim-delete"  
        And Bucket resource referencing BucketClaim resource "bucket-claim-delete" is created in ObjectStore "object-store-1"
        And BucketClaim resource "bucket-claim-delete" in namespace "namespace-1" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-delete" bucketID is not empty
        When BucketClaim resource "my-bucket-claim-delete" is deleted in namespace "namespace-1"
        Then Bucket referencing BucketClaim resource "my-bucket-claim-delete" is deleted in ObjectStore "object-store-1"

    Scenario: BucketClaim deletion with deletionPolicy set to "retain" (default)
        Given specification of custom resource "my-bucket-class-retain" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketClass
        metadata:
            name: my-bucket-class-retain
        deletionPolicy: retain
        driverName: cosi-driver  
        parameters:
            objectScaleID: ${objectScaleID}            
            objectStoreID: ${objectStoreID}
            accountSecret: ${secretName}               
        """
        And specification of custom resource "my-bucket-claim-retain" is:
        """
        apiVersion: v1
        kind: BucketClaim
        metadata:
            name: my-bucket-claim-retain
            namespace: namespace-1
        spec:                                            
            bucketClassName: my-bucket-class-retain
            protocol: S3
        """  
        And BucketClass resource is created from specification "my-bucket-class-retain"
        And BucketClaim resource is created from specification "my-bucket-claim-retain"
        And Bucket resource referencing BucketClaim resource "bucket-claim-retain" is created in ObjectStore "object-store-1"
        And BucketClaim resource "bucket-claim-retain" in namespace "namespace-1" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "bucket-claim-retain" bucketID is not empty
        When BucketClaim resource "my-bucket-claim-retain" is deleted in namespace "namespace-1"
        Then Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "object-store-1"
