//go:build integration

package main_test

import (
	"github.com/dell/cosi-driver/tests/integration/steps"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

var _ = Describe("Bucket Deletion", Serial, Label("delete"), func() {
	// Resources for scenarios
	var (
		bucketClassDelete *v1alpha1.BucketClass
		bucketClassRetain *v1alpha1.BucketClass
		bucketClaimDelete *v1alpha1.BucketClaim
		bucketClaimRetain *v1alpha1.BucketClaim
		deleteBucket      *v1alpha1.Bucket
		retainBucket      *v1alpha1.Bucket
	)

	// Background
	BeforeEach(func(ctx SpecContext) {
		// Initialize variables
		bucketClassDelete = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class-delete",
			},
			DeletionPolicy: "delete",
			DriverName:     "cosi-driver",
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		bucketClassRetain = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class-retain",
			},
			DeletionPolicy: "retain",
			DriverName:     "cosi-driver",
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		bucketClaimDelete = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-claim-delete",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class-delete",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		bucketClaimRetain = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-claim-retain",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class-retain",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		deleteBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class-delete",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "my-bucket-claim-delete", Namespace: "namespace-1"},
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		retainBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class-retain",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "my-bucket-claim-retain", Namespace: "namespace-1"},
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
		steps.CheckObjectScaleInstallation(ctx, objectscale)

		// STEP: ObjectStore "objectstore-dev" is created
		By("Checking if the ObjectStore 'objectstore-dev' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, "objectstore-dev")

		// STEP: Kubernetes namespace "driver-ns" is created
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(ctx, clientset, "driver-ns")

		// STEP: Kubernetes namespace "namespace-1" is created
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(ctx, clientset, "namespace-1")

		// STEP: COSI controller "objectstorage-controller" is installed in namespace "default"
		By("Checking if COSI controller objectstorage-controller is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		// STEP: COSI driver "cosi-driver" is installed in namespace "driver-ns"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi-driver", "driver-ns")

		DeferCleanup(func() {
			// Cleanup for background
		})
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "delete"
	It("Delets the bucket when deletionPolicy is set to 'delete'", func(ctx SpecContext) {
		// STEP: BucketClass resource is created from specification "my-bucket-class-delete"
		By("creating a BucketClass resource from specification 'my-bucket-class-delete'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassDelete)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim-delete"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-delete'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is created in ObjectStore "objectstore-dev"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is created in ObjectStore 'objectstore-dev'")
		steps.CheckBucketResourceInObjectStore(objectscale, deleteBucket)

		// STEP: BucketClaim resource "bucket-claim-delete" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-delete' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimDelete, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete'  is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, deleteBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is not empty")
		steps.CheckBucketID(ctx, bucketClient, deleteBucket)

		// STEP: BucketClaim resource "my-bucket-claim-delete" is deleted in namespace "namespace-1"
		By("deleting BucketClaim resource 'my-bucket-claim-delete' in namespace 'namespace-1'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-delete" is deleted in ObjectStore "objectstore-dev"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-delete' is deleted in ObjectStore 'objectstore-dev'")
		steps.CheckBucketDeletionInObjectStore(objectscale, deleteBucket)

		DeferCleanup(func() {
			// Cleanup for scenario: BucketClaim deletion with deletionPolicy set to "delete"
		})
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "retain"
	It("Does not delete the bucket when deletionPolicy is set to 'retain'", func(ctx SpecContext) {
		// STEP: BucketClass resource is created from specification "my-bucket-class-retain"
		By("creating a BucketClass resource from specification 'my-bucket-class-retain'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassRetain)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim-retain"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-retain'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" is created in ObjectStore "objectstore-dev"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is created in ObjectStore 'objectstore-dev'")
		steps.CheckBucketResourceInObjectStore(objectscale, retainBucket)

		// STEP: BucketClaim resource "bucket-claim-retain" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-retain' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimRetain, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-retain'  is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, retainBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketID" is not empty
		By("checking the ID of Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is not empty")
		steps.CheckBucketID(ctx, bucketClient, retainBucket)

		// STEP: BucketClaim resource "my-bucket-claim-retain" is deleted in namespace "namespace-1"
		By("deleting BucketClaim resource 'my-bucket-claim-retain' in namespace 'namespace-1'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "objectstore-dev"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-retain' is available in ObjectStore 'objectstore-dev'")
		steps.CheckBucketResourceInObjectStore(objectscale, retainBucket)

		DeferCleanup(func() {
			// Cleanup for scenario: BucketClaim deletion with deletionPolicy set to "retain"
		})
	})
})
