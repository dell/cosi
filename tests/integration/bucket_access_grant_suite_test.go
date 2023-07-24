// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

	objscl "github.com/dell/cosi-driver/pkg/provisioner/objectscale"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	"github.com/dell/cosi-driver/pkg/provisioner/policy"
	"github.com/dell/cosi-driver/tests/integration/steps"
)

var _ = Describe("Bucket Access Grant", Ordered, Label("grant", "objectscale"), func() {
	// Resources for scenarios
	var (
		myBucketClass       *v1alpha1.BucketClass
		myBucketClaim       *v1alpha1.BucketClaim
		myBucket            *v1alpha1.Bucket
		myBucketAccessClass *v1alpha1.BucketAccessClass
		myBucketAccess      *v1alpha1.BucketAccess
		validSecret         *v1.Secret
		myBucketPolicy      policy.Document
		principalUsername   string
	)

	// Background
	BeforeEach(func(ctx context.Context) {
		// Initialize variables
		myBucketClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "access-bucket-class",
			},
			DriverName:     "cosi.dellemc.com",
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		myBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "access-bucket-claim",
				Namespace: "access-namespace",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "access-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		myBucketAccessClass = &v1alpha1.BucketAccessClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccessClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "access-bucket-access-class",
			},
			DriverName:         "cosi.dellemc.com",
			AuthenticationType: v1alpha1.AuthenticationTypeKey,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		myBucketAccess = &v1alpha1.BucketAccess{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccess",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "access-bucket-access",
				Namespace: "access-namespace",
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "access-bucket-access-class",
				BucketClaimName:       "access-bucket-claim",
				CredentialsSecretName: "access-bucket-credentials",
			},
		}
		validSecret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "access-bucket-credentials",
				Namespace: "access-namespace",
			},
			Data: map[string][]byte{
				"BucketInfo": []byte(""),
			},
		}

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale)

		// STEP: ObjectStore "${objectstoreId}" is created
		By("Checking if the ObjectStore '${objectstoreId}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, ObjectstoreID)

		// STEP: Kubernetes namespace "cosi-driver" is created
		By("Checking if namespace 'cosi-driver' is created")
		steps.CreateNamespace(ctx, clientset, "cosi-driver")

		// STEP: Kubernetes namespace "access-namespace" is created
		By("Checking if namespace 'access-namespace' is created")
		steps.CreateNamespace(ctx, clientset, "access-namespace")

		// STEP: COSI controller "objectstorage-controller" is installed in namespace "default"
		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		// STEP: COSI driver "cosi-driver" is installed in namespace "cosi-driver"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'cosi-driver'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi-driver", "cosi-driver")

		// STEP: BucketClass resource is created from specification "my-bucket-class"
		By("Creating the BucketClass 'my-bucket-class' is created")
		myBucketClass = steps.CreateBucketClassResource(ctx, bucketClient, myBucketClass)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim"
		By("Creating the BucketClaim 'my-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource 'my-bucket-access-class' is created
		By("Checking if Bucket resource referencing BucketClaim resource 'my-bucket-access-class' is created")
		myBucket = steps.GetBucketResource(ctx, bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "${objectstoreName}"
		By("Checking if bucket referencing 'my-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, myBucket)

		// STEP: BucketClaim resource "my-bucket-claim" in namespace "access-namespace" status "bucketReady" is "true"
		By("Checking if BucketClaim resource 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, myBucketClaim, true)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
		By("Checking if Bucket resource referencing 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(myBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket" bucketID is not empty
		By("Checking if Bucket resource 'my-bucket' status 'bucketID' is not empty")
		steps.CheckBucketID(myBucket)

		// I need bucket name here that is generated by one of the steps above, I think. Maybe there is a better way to do this.
		resourceARN := objscl.BuildResourceString(ObjectscaleID, ObjectstoreID, myBucket.Name)
		principalUsername = objscl.BuildUsername(Namespace, myBucket.Name)
		principalARN := objscl.BuildPrincipalString(Namespace, myBucket.Name)
		myBucketPolicy = policy.Document{
			Version: "2012-10-17",
			Statement: []policy.StatementEntry{
				{
					Effect: "Allow",
					Action: []string{
						"*",
					},
					Resource: []string{
						resourceARN,
					},
					Principal: policy.PrincipalEntry{
						AWS: []string{
							principalARN,
						},
					},
				},
			},
		}
	})

	// STEP: Scenario: BucketAccess creation with KEY authorization mechanism
	It("Creates BucketAccess with KEY authorization mechanism", func(ctx context.Context) {
		// STEP: BucketAccessClass resource is created from specification "my-bucket-access-class"
		By("Creating BucketAccessClass resource 'my-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, myBucketAccessClass)

		// STEP: BucketAccess resource is created from specification "my-bucket-access"
		By("Creating BucketAccess resource 'my-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, myBucketAccess)

		// STEP: BucketAccess resource "my-bucket-access" status "accessGranted" is "true"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'access-namespace' status 'accessGranted' is 'true'")
		myBucketAccess = steps.CheckBucketAccessStatus(ctx, bucketClient, myBucketAccess, true)

		// STEP: User "user-1" in account on ObjectScale platform is created
		By("Checking if User 'user-1' in account on ObjectScale platform is created")
		steps.CheckUser(ctx, IAMClient, myBucket.Name, Namespace)

		// TODO: Change to happy policy
		// STEP: Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is created
		By("Checking if Policy for Bucket resource referencing BucketClaim resource 'my-bucket-claim' is created")
		steps.CheckPolicy(ctx, objectscale, myBucketPolicy, myBucket, Namespace)

		// TODO: Get AccountID from environment
		// STEP: BucketAccess resource "my-bucket-access" in namespace "access-namespace" status "accountID" is "${accountID}"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'access-namespace' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, myBucketAccess, principalUsername)

		// STEP: Secret "bucket-credentials-1" is created in namespace "access-namespace" and is not empty
		By("Checking if Secret 'bucket-credentials-1' in namespace 'access-namespace' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is accessible from Secret "bucket-credentials-1"
		By("Checking if Bucket resource referencing BucketClaim resource 'my-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(ctx, clientset, validSecret)

		DeferCleanup(func() {
			ctx := context.Background()
			steps.DeletePolicy(ctx, objectscale, myBucket, Namespace)
			steps.DeleteAccessKey(ctx, IAMClient, clientset, validSecret)
			steps.DeleteUser(ctx, IAMClient, myBucketAccess.Status.AccountID)
			steps.DeleteSecret(ctx, clientset, validSecret)
			steps.DeleteBucketAccessResource(ctx, bucketClient, myBucketAccess)
			steps.DeleteBucketAccessClassResource(ctx, bucketClient, myBucketAccessClass)
			steps.DeleteBucketClaimResource(ctx, bucketClient, myBucketClaim)
			steps.DeleteBucketClassResource(ctx, bucketClient, myBucketClass)
		})
	})
})
