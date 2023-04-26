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
	. "github.com/onsi/ginkgo/v2"

	"github.com/dell/cosi-driver/tests/integration/steps"
	"github.com/dell/cosi-driver/tests/integration/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

var _ = Describe("Bucket Creation", Ordered, Label("create", "objectscale"), func() {
	// Resources for scenarios
	var (
		myBucketClass      *v1alpha1.BucketClass
		validBucketClaim   *v1alpha1.BucketClaim
		invalidBucketClaim *v1alpha1.BucketClaim
		validBucket        *v1alpha1.Bucket
		myEvent            *v1.Event
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
				"ID": "${driverID}",
			},
		}
		validBucketClaim = &v1alpha1.BucketClaim{
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
		invalidBucketClaim = &v1alpha1.BucketClaim{
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
		myEvent = &v1.Event{
			Type:    v1.EventTypeWarning,
			Reason:  "MissingBucketClassName",
			Message: "BucketClassName not defined",
		}

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(ctx, clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale)

		// STEP: ObjectStore "objectstore" is created
		By("Checking if the ObjectStore 'objectstore' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, "objectstore")

		// STEP: Kubernetes namespace "driver-ns" is created
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(ctx, clientset, "driver-ns")

		// STEP: Kubernetes namespace "namespace-1" is created
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(ctx, clientset, "namespace-1")

		// STEP: COSI controller "objectstorage-controller" is installed in namespace "default"
		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		// STEP: COSI driver "cosi-driver" is installed in namespace "driver-ns"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi-driver", "driver-ns")

		// STEP: BucketClass resource is created from specification "my-bucket-class"
		By("Creating the BucketClass 'my-bucket-class' is created")
		steps.CreateBucketClassResource(ctx, bucketClient, myBucketClass)

		DeferCleanup(func() {
			// Cleanup for background
		})
	})

	// STEP: Scenario: Successfull bucket creation
	It("Successfully creates bucket", func(ctx SpecContext) {
		// STEP: BucketClaim resource is created from specification "bucket-claim-valid"
		By("creating a BucketClaim resource from specification 'bucket-claim-valid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, validBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created")
		validBucket = steps.GetBucketResource(ctx, bucketClient, validBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" is created in ObjectStore "objectstore-dev"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created in ObjectStore 'objectstore-dev'")

		steps.CheckBucketResourceInObjectStore(objectscale, validBucket)

		// STEP: BucketClaim resource "bucket-claim-valid" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-valid' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, validBucketClaim, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid'  is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, validBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is not empty")
		steps.CheckBucketID(ctx, bucketClient, validBucket)

		DeferCleanup(func() {
			steps.DeleteBucketClaimResource(ctx, bucketClient, validBucketClaim)
		})
	})

	// STEP: Scenario: Unsuccessfull bucket creation
	It("Unsuccessfully tries to create bucket", func(ctx SpecContext) {
		// STEP: BucketClaim resource is created from specification "bucket-claim-invalid"
		By("creating a BucketClaim resource from specification 'bucket-claim-invalid'")
		steps.CreateBucketClaimResource(ctx, bucketClient, invalidBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource 'bucket-claim-invalid' is created
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-invalid' is created")
		_ = steps.GetBucketResource(ctx, bucketClient, invalidBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-invalid" is not created in ObjectStore "objectstore-dev"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-invalid' is not created in ObjectStore 'objectstore-dev'")
		steps.CheckBucketNotInObjectStore(objectscale, invalidBucketClaim)

		// STEP: BucketClaim resource "bucket-claim-invalid" in namespace "namespace-1" status "bucketReady" is "false"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-invalid' in namespace 'namespace-1' is 'false'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, invalidBucketClaim, false)

		// STEP: BucketClaim events contains an error: "Cannot create Bucket: BucketClass does not exist"
		By("checking if the BucketClaim events contains an error: 'Cannot create Bucket: BucketClass does not exist'")
		steps.CheckBucketClaimEvents(ctx, clientset, invalidBucketClaim, myEvent)

		DeferCleanup(func() {
			steps.DeleteBucketClaimResource(ctx, bucketClient, invalidBucketClaim)
		})
	})
	AfterAll(func() {
		DeferCleanup(func(ctx SpecContext) {
			steps.DeleteBucketClassResource(ctx, bucketClient, myBucketClass)
			utils.DeleteReleasesAndNamespaces(ctx, clientset, map[string]string{"ns-driver": "cosi-driver"}, []string{"ns-driver"})
		})
	})
})
