//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build integration

package main_test

import (
	"context"

	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dell/cosi-driver/tests/integration/steps"
)

var _ = Describe("Bucket Deletion", Ordered, Label("delete", "objectscale"), func() {
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
				Namespace: "deletion-namespace",
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
				Namespace: "deletion-namespace",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "deletion-bucket-class-retain",
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

		// STEP: ObjectStore "${objectstoreId}" is created
		By("Checking if the ObjectStore '${objectstoreId}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, ObjectstoreID)

		// STEP: Kubernetes namespace "cosi-driver" is created
		By("Checking if namespace 'cosi-driver' is created")
		steps.CreateNamespace(ctx, clientset, "cosi-driver")

		// STEP: Kubernetes namespace "deletion-namespace" is created
		By("Checking if namespace 'deletion-namespace' is created")
		steps.CreateNamespace(ctx, clientset, "deletion-namespace")

		// STEP: COSI controller "objectstorage-controller" is installed in namespace "default"
		By("Checking if COSI controller objectstorage-controller is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		// STEP: COSI driver "cosi-driver" is installed in namespace "cosi-driver"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'cosi-driver'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi-driver", "cosi-driver")
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "delete"
	It("Deletes the bucket when deletionPolicy is set to 'delete'", func(ctx context.Context) {
		// STEP: BucketClass resource is created from specification "my-bucket-class-delete"
		By("creating a BucketClass resource from specification 'my-bucket-class-delete'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassDelete)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim-delete"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-delete'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim-delete' is created
		By("checking if Bucket resource referencing BucketClaim resource 'my-bucket-claim-delete' is created")
		deleteBucket = steps.GetBucketResource(ctx, bucketClient, bucketClaimDelete)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is created in ObjectStore "${objectstoreName}"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, deleteBucket)

		// STEP: BucketClaim resource "bucket-claim-delete" in namespace "deletion-namespace" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-delete' in namespace 'deletion-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimDelete, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete'  is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, deleteBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is not empty")
		steps.CheckBucketID(ctx, bucketClient, deleteBucket)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "${objectstoreName}"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-retain' is available in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, deleteBucket)

		// STEP: BucketClaim resource "my-bucket-claim-delete" is deleted in namespace "deletion-namespace"
		By("deleting BucketClaim resource 'my-bucket-claim-delete' in namespace 'deletion-namespace'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimDelete)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-delete" is deleted in ObjectStore "${objectstoreName}"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-delete' is deleted in ObjectStore '${objectstoreName}'")
		steps.CheckBucketDeletionInObjectStore(ctx, objectscale, Namespace, deleteBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClassResource(ctx, bucketClient, bucketClassDelete)
		})
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "retain"
	It("Does not delete the bucket when deletionPolicy is set to 'retain'", func(ctx context.Context) {
		// STEP: BucketClass resource is created from specification "my-bucket-class-retain"
		By("creating a BucketClass resource from specification 'my-bucket-class-retain'")
		steps.CreateBucketClassResource(ctx, bucketClient, bucketClassRetain)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim-retain"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-retain'")
		steps.CreateBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		// STEP: Bucket resource referencing BucketClaim resource 'my-bucket-claim-retain' is created"
		By("checking if Bucket resource referencing BucketClaim resource 'my-bucket-claim-retain' is created")
		retainBucket = steps.GetBucketResource(ctx, bucketClient, bucketClaimRetain)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" is created in ObjectStore "${objectstoreName}"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, retainBucket)

		// STEP: BucketClaim resource "bucket-claim-retain" in namespace "deletion-namespace" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-retain' in namespace 'deletion-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, bucketClaimRetain, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-retain'  is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, retainBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketID" is not empty
		By("checking the ID of Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is not empty")
		steps.CheckBucketID(ctx, bucketClient, retainBucket)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "${objectstoreId}"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-retain' is available in ObjectStore '${objectstoreId}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, retainBucket)

		// STEP: BucketClaim resource "my-bucket-claim-retain" is deleted in namespace "deletion-namespace"
		By("deleting BucketClaim resource 'my-bucket-claim-retain' in namespace 'deletion-namespace'")
		steps.DeleteBucketClaimResource(ctx, bucketClient, bucketClaimRetain)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "${objectstoreId}"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-retain' is available in ObjectStore '${objectstoreId}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, retainBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClassResource(ctx, bucketClient, bucketClassRetain)
		})
	})
})
