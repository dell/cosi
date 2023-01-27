//go:build integration

package main_test

import (
	"github.com/dell/cosi-driver/tests/integration/steps"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Bucket Creation", Label("create"), func() {
	// Resources for scenarios
	var (
		myBucketClass      *v1alpha1.BucketClass
		bucketClaimValid   *v1alpha1.BucketClaim
		bucketClaimInvalid *v1alpha1.BucketClaim
		validBucket        *v1alpha1.Bucket
	)

	// Background
	BeforeEach(func(ctx SpecContext) {
		// Initialize variables
		myBucketClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class",
			},
			DeletionPolicy: "delete",
			DriverName:     "cosi-driver",
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		bucketClaimValid = &v1alpha1.BucketClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bucket-claim-valid",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		bucketClaimInvalid = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bucket-claim-invalid",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "bucket-class-invalid",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		validBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "bucket-claim-valid", Namespace: "namespace-1"},
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(ctx, clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, clientset)

		// STEP: ObjectStore "object-store-1" is created
		By("Checking if the ObjectStore 'object-store-1' is created")
		steps.CreateObjectStore(ctx, objectscale, "object-store-1")

		// STEP: Kubernetes namespace "driver-ns" is created
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(ctx, clientset, "driver-ns")

		// STEP: Kubernetes namespace "namespace-1" is created
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(ctx, clientset, "namespace-1")

		// STEP: COSI controller "cosi-controller" is installed in namespace "driver-ns"
		By("Checking if COSI controller 'cosi-controller' is installed in namespace 'driver-ns'")
		steps.CheckCOSIControllerInstallation(clientset, "cosi-controller", "driver-ns")

		// STEP: COSI driver "cosi-driver" is installed in namespace "driver-ns"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(clientset, "cosi-driver", "driver-ns")

		// STEP: BucketClass resource is created from specification "my-bucket-class"
		By("Creating the BucketClass 'my-bucket-class' is created")
		steps.CreateBucketClassResource(bucketClient, myBucketClass)

		DeferCleanup(func() {
			// Cleanup for background
		})
	})

	// STEP: Scenario: Successfull bucket creation
	It("Successfully creates bucket", func() {
		// STEP: BucketClaim resource is created from specification "bucket-claim-valid"
		By("creating a BucketClaim resource from specification 'bucket-claim-valid'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimValid)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" is created in ObjectStore "object-store-1"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, validBucket)

		// STEP: BucketClaim resource "bucket-claim-valid" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-valid' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimValid)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid'  is 'true'")
		steps.CheckBucketStatus(bucketClient, validBucket)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is not empty")
		steps.CheckBucketID(bucketClient, validBucket)

		DeferCleanup(func() {
			// Cleanup for scenario: Successfull bucket creation
		})
	})

	// STEP: Scenario: Unsuccessfull bucket creation
	It("Unsuccessfully tries to create bucket", func() {
		// STEP: BucketClaim resource is created from specification "bucket-claim-invalid"
		By("creating a BucketClaim resource from specification 'bucket-claim-invalid'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimInvalid)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-invalid" is not created in ObjectStore "object-store-1"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-invalid' is not created in ObjectStore 'object-store-1'")
		steps.CheckBucketNotInObjectStore(objectscale, bucketClaimInvalid)

		// STEP: BucketClaim resource "bucket-claim-invalid" in namespace "namespace-1" status "bucketReady" is "false"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-invalid' in namespace 'namespace-1' is 'false'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimInvalid)

		// STEP: BucketClaim events contains an error: "Cannot create Bucket: BucketClass does not exist"
		By("checking if the BucketClaim events contains an error: 'Cannot create Bucket: BucketClass does not exist'")
		steps.CheckBucketClaimEvents(clientset, bucketClaimInvalid)

		DeferCleanup(func() {
			// Cleanup for scenario: Unsuccessfull bucket creation
		})
	})

})
