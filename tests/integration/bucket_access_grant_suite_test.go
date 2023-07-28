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

	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/cosi/tests/integration/steps"
)

var _ = Describe("Bucket Access Grant", Ordered, Label("grant", "objectscale"), func() {
	// Resources for scenarios
	var (
		grantBucketClass       *v1alpha1.BucketClass
		grantBucketClaim       *v1alpha1.BucketClaim
		grantBucket            *v1alpha1.Bucket
		grantBucketAccessClass *v1alpha1.BucketAccessClass
		grantBucketAccess      *v1alpha1.BucketAccess
		validSecret            *v1.Secret
		grantBucketPolicy      policy.Document
		principalUsername      string
	)

	// Background
	BeforeEach(func(ctx context.Context) {
		// Initialize variables
		grantBucketClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "grant-bucket-class",
			},
			DriverName:     "cosi.dellemc.com",
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		grantBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "grant-bucket-claim",
				Namespace: "access-grant-namespace",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "grant-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		grantBucketAccessClass = &v1alpha1.BucketAccessClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccessClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "grant-bucket-access-class",
			},
			DriverName:         "cosi.dellemc.com",
			AuthenticationType: v1alpha1.AuthenticationTypeKey,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		grantBucketAccess = &v1alpha1.BucketAccess{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccess",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "grant-bucket-access",
				Namespace: "access-grant-namespace",
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "grant-bucket-access-class",
				BucketClaimName:       "grant-bucket-claim",
				CredentialsSecretName: "grant-bucket-credentials",
			},
		}
		validSecret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "grant-bucket-credentials",
				Namespace: "access-grant-namespace",
			},
			Data: map[string][]byte{
				"BucketInfo": []byte(`{"metadata":{"name":""},"spec":{"bucketName":"","authenticationType":"","secretS3":{"endpoint":"","region":"","accessKeyID":"","accessSecretKey":""},"protocols":[]}}`),
			},
		}

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale, Namespace)

		// STEP: ObjectStore "${objectstoreId}" is created
		By("Checking if the ObjectStore '${objectstoreId}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, ObjectstoreID)

		// STEP: Kubernetes namespace "cosi-driver" is created
		By("Checking if namespace 'cosi-driver' is created")
		steps.CreateNamespace(ctx, clientset, "cosi-driver")

		// STEP: Kubernetes namespace "access-grant-namespace" is created
		By("Checking if namespace 'access-grant-namespace' is created")
		steps.CreateNamespace(ctx, clientset, "access-grant-namespace")

		// STEP: COSI controller "objectstorage-controller" is installed in namespace "default"
		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		// STEP: COSI driver "cosi-driver" is installed in namespace "cosi-driver"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'cosi-driver'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, "cosi-driver", "cosi-driver")

		// STEP: BucketClass resource is created from specification "grant-bucket-class"
		By("Creating the BucketClass 'grant-bucket-class' is created")
		grantBucketClass = steps.CreateBucketClassResource(ctx, bucketClient, grantBucketClass)

		// STEP: BucketClaim resource is created from specification "grant-bucket-claim"
		By("Creating the BucketClaim 'grant-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, grantBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource 'grant-bucket-access-class' is created
		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-access-class' is created")
		grantBucket = steps.GetBucketResource(ctx, bucketClient, grantBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "grant-bucket-claim" is created in ObjectStore "${objectstoreName}"
		By("Checking if bucket referencing 'grant-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, grantBucket)

		// STEP: BucketClaim resource "grant-bucket-claim" in namespace "access-grant-namespace" status "bucketReady" is "true"
		By("Checking if BucketClaim resource 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, grantBucketClaim, true)

		// STEP: Bucket resource referencing BucketClaim resource "grant-bucket-claim" status "bucketReady" is "true"
		By("Checking if Bucket resource referencing 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(grantBucket, true)

		// STEP: Bucket resource referencing BucketClaim resource "grant-bucket" bucketID is not empty
		By("Checking if Bucket resource 'grant-bucket' status 'bucketID' is not empty")
		steps.CheckBucketID(grantBucket)

		// I need bucket name here that is generated by one of the steps above, I think. Maybe there is a better way to do this.
		resourceARN := objscl.BuildResourceString(ObjectscaleID, ObjectstoreID, grantBucket.Name)
		principalUsername = objscl.BuildUsername(Namespace, grantBucket.Name)
		principalARN := objscl.BuildPrincipalString(Namespace, grantBucket.Name)
		grantBucketPolicy = policy.Document{
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
		// STEP: BucketAccessClass resource is created from specification "grant-bucket-access-class"
		By("Creating BucketAccessClass resource 'grant-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, grantBucketAccessClass)

		// STEP: BucketAccess resource is created from specification "grant-bucket-access"
		By("Creating BucketAccess resource 'grant-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, grantBucketAccess)

		// STEP: BucketAccess resource "grant-bucket-access" status "accessGranted" is "true"
		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accessGranted' is 'true'")
		grantBucketAccess = steps.CheckBucketAccessStatus(ctx, bucketClient, grantBucketAccess, true)

		// STEP: User "user-1" in account on ObjectScale platform is created
		By("Checking if User 'user-1' in account on ObjectScale platform is created")
		steps.CheckUser(ctx, IAMClient, grantBucket.Name, Namespace)

		// STEP: Policy "${policy}" for Bucket resource referencing BucketClaim resource "grant-bucket-claim" on ObjectScale platform is created
		By("Checking if Policy for Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is created")
		steps.CheckPolicy(ctx, objectscale, grantBucketPolicy, grantBucket, Namespace)

		// STEP: BucketAccess resource "grant-bucket-access" in namespace "access-grant-namespace" status "accountID" is "${accountID}"
		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, grantBucketAccess, principalUsername)

		// STEP: Secret "bucket-credentials-1" is created in namespace "access-grant-namespace" and is not empty
		By("Checking if Secret 'bucket-credentials-1' in namespace 'access-grant-namespace' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is accessible from Secret "bucket-credentials-1"
		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(ctx, clientset, validSecret)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketAccessResource(ctx, bucketClient, grantBucketAccess)
			steps.DeleteBucketAccessClassResource(ctx, bucketClient, grantBucketAccessClass)

			revokedPolicy := policy.Document{
				Version: "2012-10-17",
				Statement: []policy.StatementEntry{
					{
						Effect: "Allow",
						Action: []string{
							"*",
						},
						Resource: []string{},
						Principal: policy.PrincipalEntry{
							AWS: []string{},
						},
					},
				},
			}

			// We should wait until access to policy is revoked, to not cause errors.
			steps.CheckPolicy(ctx, objectscale, revokedPolicy, grantBucket, Namespace)

			steps.DeleteBucketClaimResource(ctx, bucketClient, grantBucketClaim)
			steps.DeleteBucketClassResource(ctx, bucketClient, grantBucketClass)
		})
	})
})
