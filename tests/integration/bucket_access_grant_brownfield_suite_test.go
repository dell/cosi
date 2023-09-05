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

	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"
	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/cosi/tests/integration/steps"
)

var _ = Describe("Bucket Access Grant for Brownfield Bucket", Ordered, Label("grant", "brownfield", "objectscale"), func() {
	// Resources for scenarios
	const (
		namespace string = "access-grant-namespace-brownfield"
	)

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
		grantBucket = &v1alpha1.Bucket{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Bucket",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-brownfield-bucket",
			},
			Spec: v1alpha1.BucketSpec{
				BucketClaim:      &v1.ObjectReference{},
				BucketClassName:  "brownfield-grant-bucket-class",
				DriverName:       "cosi.dellemc.com",
				DeletionPolicy:   "Retain",
				ExistingBucketID: DriverID + "-my-brownfield-bucket",
				Parameters: map[string]string{
					"id": DriverID,
				},
				Protocols: []v1alpha1.Protocol{v1alpha1.ProtocolS3},
			},
		}
		grantBucketClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "brownfield-grant-bucket-class",
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
				Name:      "brownfield-grant-bucket-claim",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName:    "brownfield-grant-bucket-class",
				ExistingBucketName: "my-brownfield-bucket",
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
				Name: "brownfield-grant-bucket-access-class",
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
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "brownfield-grant-bucket-access-class",
				BucketClaimName:       "brownfield-grant-bucket-claim",
				CredentialsSecretName: "brownfield-grant-bucket-credentials",
			},
		}
		validSecret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "brownfield-grant-bucket-credentials",
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"BucketInfo": []byte(`{"metadata":{"name":""},"spec":{"bucketName":"","authenticationType":"","secretS3":{"endpoint":"","region":"","accessKeyID":"","accessSecretKey":""},"protocols":[]}}`),
			},
		}

		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, objectscale, Namespace)

		By("Checking if the ObjectStore '${objectstoreId}' is created")
		steps.CheckObjectStoreExists(ctx, objectscale, ObjectstoreID)

		By("Checking if namespace 'cosi-test-ns' is created")
		steps.CreateNamespace(ctx, clientset, DriverNamespace)

		By("Checking if namespace 'access-grant-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)

		By("Creating the BucketClass 'grant-bucket-class' is created")
		grantBucketClass = steps.CreateBucketClassResource(ctx, bucketClient, grantBucketClass)

		By("Creating bucket on the Objectscale platform")
		steps.CreateBucket(ctx, objectscale, Namespace, grantBucket)

		By("Creating Bucket")
		steps.CreateBucketResource(ctx, bucketClient, grantBucket)

		By("Creating the BucketClaim 'grant-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, grantBucketClaim)

		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-access-class' is created")
		grantBucket = steps.GetBucketResource(ctx, bucketClient, grantBucketClaim)

		By("Checking if bucket referencing 'grant-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, grantBucket)

		By("Checking if BucketClaim resource 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, grantBucketClaim, true)

		By("Checking if Bucket resource referencing 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(grantBucket, true)

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

	It("Creates BucketAccess with KEY authorization mechanism", func(ctx context.Context) {
		By("Creating BucketAccessClass resource 'grant-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, grantBucketAccessClass)

		By("Creating BucketAccess resource 'grant-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, grantBucketAccess)

		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accessGranted' is 'true'")
		grantBucketAccess = steps.CheckBucketAccessStatus(ctx, bucketClient, grantBucketAccess, true)

		By("Checking if User 'user-1' in account on ObjectScale platform is created")
		steps.CheckUser(ctx, IAMClient, grantBucket.Name, Namespace)

		By("Checking if Policy for Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is created")
		steps.CheckPolicy(ctx, objectscale, grantBucketPolicy, grantBucket, Namespace)

		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, grantBucketAccess, principalUsername)

		By("Checking if Secret 'bucket-credentials-1' in namespace 'access-grant-namespace' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret)

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
			steps.DeleteBucket(ctx, objectscale, Namespace, grantBucket)
		})
	})
})
