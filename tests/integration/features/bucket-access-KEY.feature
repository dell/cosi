@component_COSI

Feature: BucketAccess creation in KEY flow on ObjectScale platform

    As an ObjectScale platform user
    I want to add BucketAccess via KEY authentication, which is a request access to a Bucket for particular account
    so that access credentials for a Bucket are created and unique identifier for the account (accountID) is returned

    Background:
        Given Kubernetes cluster is up and running
        And ObjectScale platform is installed on the cluster
        And ObjectStore "${objectstoreName}" is created
        And Kubernetes namespace "driver-ns" is created
        And Kubernetes namespace "namespace-1" is created
        And COSI controller "objectstorage-controller" is installed in namespace "default"
        And COSI driver "cosi-driver" is installed in namespace "driver-ns"
        And specification of custom resource "my-bucket-class" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketClass
        metadata:
            name: my-bucket-class
        deletionPolicy: delete
        driverName: cosi-driver
        parameters:
            ID: ${driverID}
        """
        And specification of custom resource "my-bucket-claim" is:
        """
        apiVersion: v1
        kind: BucketClaim
        metadata:
            name: my-bucket-claim
            namespace: namespace-1
        spec:
            bucketClassName: my-bucket-class
            protocol: S3
        """
        And BucketClass resource is created from specification "my-bucket-class"
        And BucketClaim resource is created from specification "my-bucket-claim"
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" is created
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "${objectstoreName}"
        And BucketClaim resource "my-bucket-claim" in namespace "namespace-1" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" bucketID is not empty

    Scenario: BucketAccess creation with KEY authorization mechanism 
        And specification of custom resource "my-bucket-access-class" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketAccessClass
        metadata:
            name: my-bucket-access-class
        driverName: cosi-driver
        authenticationType: KEY
        parameters:
            ID: ${driverID}
        """
        And specification of custom resource "my-bucket-access" is:
        """
        apiVersion: v1
        kind: BucketAccess
        metadata:
            name: my-bucket-access
            namespace: namespace-1
        spec:
            bucketAccessClassName: my-bucket-access-class
            bucketClaimName: my-bucket-claim
            credentialsSecretName: bucket-credentials-1
        """
        When BucketAccessClass resource is created from specification "my-bucket-access-class"
        And BucketAccess resource is created from specification "my-bucket-access"
        Then BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accessGranted" is "true"
        And User "${user}" in account on ObjectScale platform is created
        And Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is created
        And BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accountID" is "${accountID}"
        And Secret "bucket-credentials-1" is created in namespace "namespace-1" and is not empty
        And Bucket resource referencing BucketClaim resource "bucket-claim-delete" is accessible from Secret "bucket-credentials-1"
