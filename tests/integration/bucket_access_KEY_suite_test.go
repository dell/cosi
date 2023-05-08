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

var _ = Describe("Bucket Access KEY", Ordered, Label("key-flow", "objectscale"), func() {
	// Resources for scenarios
	var (
		myBucketClass       *v1alpha1.BucketClass
		myBucketClaim       *v1alpha1.BucketClaim
		myBucket            *v1alpha1.Bucket
		myBucketAccessClass *v1alpha1.BucketAccessClass
		myBucketAccess      *v1alpha1.BucketAccess
		validSecret         *v1.Secret
	)

	// Background
	BeforeEach(func(ctx SpecContext) {
		// Initialize variables
		myBucketClass = &v1alpha1.BucketClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class",
			},
			DriverName:     "cosi-driver",
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			Parameters: map[string]string{
				"id": driverID,
			},
		}
		myBucketClaim = &v1alpha1.BucketClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-claim",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		myBucketAccessClass = &v1alpha1.BucketAccessClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-access-class",
			},
			DriverName:         "cosi-driver",
			AuthenticationType: v1alpha1.AuthenticationTypeKey,
			Parameters: map[string]string{
				"id": driverID,
			},
		}
		myBucketAccess = &v1alpha1.BucketAccess{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-access",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "my-bucket-access-class",
				BucketClaimName:       "my-bucket-claim",
				CredentialsSecretName: "bucket-credentials-1",
			},
		}
		validSecret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "valid-secret-1",
				Namespace: "namespace-1",
			},
			Data: map[string][]byte{
				"this": []byte("is template for data"), // FIXME: when we know exact format of the secret
			},
		}

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(ctx, clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale)

		// STEP: ObjectStore "${objectstoreName}" is created
		By("Checking if the ObjectStore '${objectstoreName}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, objectstoreName)

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

		// STEP: BucketClaim resource is created from specification "my-bucket-claim"
		By("Creating the BucketClaim 'my-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource 'my-bucket-access-class' is created
		By("Checking if Bucket resource referencing BucketClaim resource 'my-bucket-access-class' is created")
		myBucket = steps.GetBucketResource(ctx, bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "${objectstoreName}"
		By("Checking if bucket referencing 'my-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(objectscale, namespace, myBucket)

		// STEP: BucketClaim resource "my-bucket-claim" in namespace "namespace-1" status "bucketReady" is "true"
		By("Checking if BucketClaim resource 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, myBucketClaim, true)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
		By("Checking if Bucket resource referencing 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(ctx, bucketClient, myBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket" bucketID is not empty
		By("Checking if Bucket resource 'my-bucket' status 'bucketID' is not empty")
		steps.CheckBucketID(ctx, bucketClient, myBucket)

		DeferCleanup(func() {
			// Cleanup for background
		})
	})

	// STEP: Scenario: BucketAccess creation with KEY authorization mechanism
	It("Creates BucketAccess with KEY authorization mechanism", func(ctx SpecContext) {
		// STEP: BucketAccessClass resource is created from specification "my-bucket-access-class"
		By("Creating BucketAccessClass resource 'my-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, myBucketAccessClass)

		// STEP: BucketAccess resource is created from specification "my-bucket-access"
		By("Creating BucketAccess resource 'my-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, myBucketAccess)

		// STEP: BucketAccess resource "my-bucket-access" status "accessGranted" is "true"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accessGranted' is 'true'")
		steps.CheckBucketAccessStatus(ctx, bucketClient, myBucketAccess, true)

		// STEP: User "${user}" in account on ObjectScale platform is created
		By("Checking if User '${user}' in account on ObjectScale platform is created")
		steps.CheckUser(ctx, iamClient, "${user}")

		// STEP: Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is created
		By("Checking if Policy '${policy}' for Bucket resource referencing BucketClaim resource 'my-bucket-claim' is created")
		steps.CheckPolicy(objectscale, "${policy}", myBucket)

		// STEP: BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accountID" is "${accountID}"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, myBucketAccess, "${accountID}")

		// STEP: Secret "bucket-credentials-1" is created in namespace "namespace-1" and is not empty
		By("Checking if Secret 'bucket-credentials-1' in namespace 'namespace-1' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is accessible from Secret "bucket-credentials-1"
		By("Checking if Bucket resource referencing BucketClaim resource 'my-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(objectscale, myBucket, "bucket-credentials-1")

		DeferCleanup(func() {
			steps.DeleteBucketAccessResource(ctx, bucketClient, myBucketAccess)
			steps.DeletePolicy(objectscale, myBucket)
			steps.DeleteUser(ctx, iamClient, "${user}")
			steps.DeleteBucketAccessClassResource(ctx, bucketClient, myBucketAccessClass)
		})
	})
	AfterAll(func() {
		DeferCleanup(func(ctx SpecContext) {
			steps.DeleteBucketClassResource(ctx, bucketClient, myBucketClass)
			steps.DeleteBucketClaimResource(ctx, bucketClient, myBucketClaim)
			utils.DeleteReleasesAndNamespaces(ctx, clientset, map[string]string{"ns-driver": "cosi-driver"}, []string{"ns-driver"})
		})
	})
})
