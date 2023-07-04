@component_COSI

Feature: BucketAccess deletion on ObjectScale platform

    As an ObjectScale platform user
    I want to delete BucketAccess
    so that access for a Bucket is deleted for particular account

    Background:
        Given Kubernetes cluster is up and running
        And ObjectScale platform is installed on the cluster
        And ObjectStore "${objectstoreId}" is created
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
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "${objectstoreId}"
        And BucketClaim resource "my-bucket-claim" in namespace "namespace-1" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
        And Bucket resource referencing BucketClaim resource "my-bucket-claim" bucketID is not empty
        And specification of custom resource "my-bucket-access-class" is:
        """
        apiVersion: storage.k8s.io/v1
        kind: BucketAccessClass
        metadata:
            name: my-bucket-access-class
        driverName: cosi.dellemc.com
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
        And BucketAccessClass resource is created from specification "my-bucket-access-class"
        And BucketAccess resource is created from specification "my-bucket-access"
        And BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accessGranted" is "true"
        And User "${user}" in account on ObjectScale platform is created
        And Policy "${policy}" on ObjectScale platform is created
        And BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accountID" is "${accountID}"
        And Secret "bucket-credentials-1" is created in namespace "namespace-1" and is not empty

    Scenario: Revoke access to bucket
        When BucketAccess resource "my-bucket-access" in namespace "namespace-1" is deleted
        And Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is deleted
        Then User "${user}" in account on ObjectScale platform is deleted


