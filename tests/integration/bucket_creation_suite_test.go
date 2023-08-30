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

	"github.com/dell/cosi/tests/integration/steps"
)

var _ = Describe("Bucket Creation", Ordered, Label("create", "objectscale"), func() {
	// Resources for scenarios
	var (
		createClass        *v1alpha1.BucketClass
		validBucketClaim   *v1alpha1.BucketClaim
		invalidBucketClaim *v1alpha1.BucketClaim
		validBucket        *v1alpha1.Bucket
		// TODO: waiting for event PR merge to sidecar
		// myEvent            *v1.Event
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
				Namespace: "creation-namespace",
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
				Namespace: "creation-namespace",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "creation-bucket-class-invalid",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		// TODO: waiting for event PR merge to sidecar
		// myEvent = &v1.Event{
		// 	Type:    v1.EventTypeWarning,
		// 	Reason:  "MissingBucketClassName",
		// 	Message: "BucketClassName not defined",
		// }

		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale, Namespace)

		By("Checking if the ObjectStore '${objectstoreId}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, ObjectstoreID)

		By("Checking if namespace 'cos-test-ns' is created")
		steps.CreateNamespace(ctx, clientset, "cosi-test-ns")

		By("Checking if namespace 'creation-namespace' is created")
		steps.CreateNamespace(ctx, clientset, "creation-namespace")

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi", "cosi-test-ns")

		By("Creating the BucketClass 'create-bucket-class'")
		steps.CreateBucketClassResource(ctx, bucketClient, createClass)
	})

	It("Successfully creates bucket", func(ctx context.Context) {
		By("creating a BucketClaim resource from specification 'bucket-claim-valid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, validBucketClaim)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created")
		validBucket = steps.GetBucketResource(ctx, bucketClient, validBucketClaim)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, validBucket)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-valid' in namespace 'creation-namespace' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, validBucketClaim, true)

		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid'  is 'true'")
		steps.CheckBucketStatus(validBucket, true)

		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is not empty")
		steps.CheckBucketID(validBucket)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClaimResource(ctx, bucketClient, validBucketClaim)
			// steps.DeleteBucket(ctx, objectscale, Namespace, validBucket)
		})
	})

	It("Unsuccessfully tries to create bucket", func(ctx context.Context) {
		By("creating a BucketClaim resource from specification 'bucket-claim-invalid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, invalidBucketClaim)

		By("checking if Bucket status in BucketClaim resource is empty")
		steps.CheckBucketStatusEmpty(ctx, bucketClient, invalidBucketClaim)

		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-invalid' is not created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketNotInObjectStore(ctx, objectscale, invalidBucketClaim)

		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-invalid' in namespace 'creation-namespace' is 'false'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, invalidBucketClaim, false)

		// NOTE: commented for now until changes introduced to provisioner sidecar
		// By("checking if the BucketClaim events contains an error: 'Cannot create Bucket: BucketClass does not exist'")
		// steps.CheckBucketClaimEvents(ctx, clientset, invalidBucketClaim, myEvent)

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
