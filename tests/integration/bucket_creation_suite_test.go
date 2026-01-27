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

var _ = Describe("Bucket Creation", Ordered, Label("create", "objectscale"), func() {
	// Resources for scenarios
	const (
		namespace string = "creation-namespace"
	)
	var (
		createClass        *v1alpha1.BucketClass
		validBucketClaim   *v1alpha1.BucketClaim
		invalidBucketClaim *v1alpha1.BucketClaim
		validBucket        *v1alpha1.Bucket
	)

	// Background
	BeforeEach(func(ctx context.Context) {
		// Initialize variables
		createClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "creation-bucket-class",
			},
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			DriverName:     "cosi.dellemc.com",
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		validBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "creation-bucket-claim-valid",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "creation-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}

		invalidBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "creation-bucket-claim-invalid",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "creation-bucket-class-invalid",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}

		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, mgmtClient, Namespace)

		By("Checking if namespace 'cos-test-ns' is created")
		steps.CreateNamespace(ctx, clientset, DriverNamespace)

		By("Checking if namespace 'creation-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "container-object-storage-controller", "container-object-storage-system")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)

		By("Creating the BucketClass 'create-bucket-class'")
		steps.CreateBucketClassResource(ctx, bucketClient, createClass)
	})

	It("Successfully creates bucket", func(ctx context.Context) {
		By("creating a BucketClaim resource from specification 'bucket-claim-valid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, validBucketClaim)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created")
		validBucket = steps.GetBucketResource(ctx, bucketClient, validBucketClaim)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, validBucket)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-valid' in namespace 'creation-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, validBucketClaim, true)

		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid'  is 'true'")
		steps.CheckBucketStatus(validBucket, true)

		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is not empty")
		steps.CheckBucketID(validBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClaimResource(ctx, bucketClient, validBucketClaim)
		})
	})

	It("Unsuccessfully tries to create bucket", func(ctx context.Context) {
		By("creating a BucketClaim resource from specification 'bucket-claim-invalid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, invalidBucketClaim)

		By("checking if Bucket status in BucketClaim resource is empty")
		steps.CheckBucketStatusEmpty(ctx, bucketClient, invalidBucketClaim)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-invalid' in namespace 'creation-namespace' is 'false'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, invalidBucketClaim, false)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClaimResource(ctx, bucketClient, invalidBucketClaim)
		})
	})

	AfterEach(func() {
		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClassResource(ctx, bucketClient, createClass)
		})
	})
})
