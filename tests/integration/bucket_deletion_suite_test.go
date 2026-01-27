// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

//go:build integration

package main_test

import (
	"context"

	"sigs.k8s.io/container-object-storage-interface/client/apis/objectstorage/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dell/cosi/tests/integration/steps"
)

var _ = Describe("Bucket Deletion", Ordered, Label("delete", "objectscale"), func() {
	// Resources for scenarios
	const (
		namespace string = "deletion-namespace"
	)
	var (
		bucketClassDelete *v1alpha1.BucketClass
		bucketClassRetain *v1alpha1.BucketClass
		bucketClaimDelete *v1alpha1.BucketClaim
		bucketClaimRetain *v1alpha1.BucketClaim
		deleteBucket      *v1alpha1.Bucket
		retainBucket      *v1alpha1.Bucket
	)

	// Background
	BeforeEach(func(ctx context.Context) {
		// Initialize variables
		bucketClassDelete = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "deletion-bucket-class-delete",
			},
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			DriverName:     "cosi.dellemc.com",
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		bucketClassRetain = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "deletion-bucket-class-retain",
			},
			DeletionPolicy: v1alpha1.DeletionPolicyRetain,
			DriverName:     "cosi.dellemc.com",
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		bucketClaimDelete = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deletion-bucket-claim-delete",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "deletion-bucket-class-delete",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		bucketClaimRetain = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "deletion-bucket-claim-retain",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "deletion-bucket-class-retain",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}

		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, mgmtClient, Namespace)

		By("Checking if namespace 'cosi-test-ns' is created")
		steps.CreateNamespace(ctx, clientset, DriverNamespace)

		By("Checking if namespace 'deletion-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller objectstorage-controller is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "container-object-storage-controller", "container-object-storage-system")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)
	})

	It("Deletes the bucket when deletionPolicy is set to 'delete'", func(ctx context.Context) {
		By("creating a BucketClass resource from specification 'delete-bucket-class-delete'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassDelete)

		By("creating a BucketClaim resource from specification 'delete-bucket-claim-delete'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		By("checking if Bucket resource referencing BucketClaim resource 'delete-bucket-claim-delete' is created")
		deleteBucket = steps.GetBucketResource(ctx, bucketClient, bucketClaimDelete)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, deleteBucket)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-delete' in namespace 'deletion-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimDelete, true)

		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete'  is 'true'")
		steps.CheckBucketStatus(deleteBucket, true)

		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is not empty")
		steps.CheckBucketID(deleteBucket)

		By("checking if Bucket referencing BucketClaim resource 'delete-bucket-claim-retain' is available in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, deleteBucket)

		By("deleting BucketClaim resource 'delete-bucket-claim-delete' in namespace 'deletion-namespace'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		By("checking if Bucket referencing BucketClaim resource 'delete-bucket-claim-delete' is deleted in ObjectStore '${objectstoreName}'")
		steps.CheckBucketDeletionInObjectStore(ctx, mgmtClient, Namespace, deleteBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClassResource(ctx, bucketClient, bucketClassDelete)
		})
	})

	It("Does not delete the bucket when deletionPolicy is set to 'retain'", func(ctx context.Context) {
		By("creating a BucketClass resource from specification 'delete-bucket-class-retain'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassRetain)

		By("creating a BucketClaim resource from specification 'delete-bucket-claim-retain'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		By("checking if Bucket resource referencing BucketClaim resource 'delete-bucket-claim-retain' is created")
		retainBucket = steps.GetBucketResource(ctx, bucketClient, bucketClaimRetain)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, retainBucket)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-retain' in namespace 'deletion-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimRetain, true)

		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-retain'  is 'true'")
		steps.CheckBucketStatus(retainBucket, true)

		By("checking the ID of Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is not empty")
		steps.CheckBucketID(retainBucket)

		By("checking if Bucket referencing BucketClaim resource 'delete-bucket-claim-retain' is available in ObjectStore '${objectstoreId}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, retainBucket)

		By("deleting BucketClaim resource 'delete-bucket-claim-retain' in namespace 'deletion-namespace'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		By("checking if Bucket referencing BucketClaim resource 'delete-bucket-claim-retain' is available in ObjectStore '${objectstoreId}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, retainBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClassResource(ctx, bucketClient, bucketClassRetain)
			steps.DeleteBucket(ctx, mgmtClient, Namespace, retainBucket)
		})
	})
})
