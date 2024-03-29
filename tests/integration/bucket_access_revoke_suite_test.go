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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"
	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/cosi/tests/integration/steps"
)

var _ = Describe("Bucket Access Revoke", Ordered, Label("revoke", "objectscale"), func() {
	// Resources for scenarios
	const (
		namespace string = "access-revoke-namespace"
	)

	var (
		revokeBucketClass       *v1alpha1.BucketClass
		revokeBucketClaim       *v1alpha1.BucketClaim
		revokeBucket            *v1alpha1.Bucket
		revokeBucketAccessClass *v1alpha1.BucketAccessClass
		revokeBucketAccess      *v1alpha1.BucketAccess
		initialPolicy           policy.Document
		finalPolicy             policy.Document
		validSecret             *v1.Secret
		principalUsername       string
	)

	// Background
	BeforeEach(func(ctx context.Context) {
		// Initialize variables
		revokeBucketClass = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "revoke-bucket-class",
			},
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			DriverName:     "cosi.dellemc.com",
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		revokeBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "revoke-bucket-claim",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "revoke-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		revokeBucketAccessClass = &v1alpha1.BucketAccessClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccessClass",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "revoke-bucket-access-class",
			},
			DriverName:         "cosi.dellemc.com",
			AuthenticationType: v1alpha1.AuthenticationTypeKey,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		revokeBucketAccess = &v1alpha1.BucketAccess{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccess",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "revoke-bucket-access",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "revoke-bucket-access-class",
				BucketClaimName:       "revoke-bucket-claim",
				CredentialsSecretName: "revoke-bucket-credentials",
			},
		}
		validSecret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "revoke-bucket-credentials",
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

		By("Checking if namespace 'access-revoke-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "objectstorage-controller", "default")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)

		By("Creating the BucketClass 'revoke-bucket-class'")
		steps.CreateBucketClassResource(ctx, bucketClient, revokeBucketClass)

		By("Creating the BucketClaim 'revoke-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, revokeBucketClaim)

		By("Checking if Bucket resource referencing BucketClaim resource 'revoke-bucket-claim' is created")
		revokeBucket = steps.GetBucketResource(ctx, bucketClient, revokeBucketClaim)

		By("Checking if the Bucket referencing 'revoke-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, objectscale, Namespace, revokeBucket)

		By("Checking if the BucketClaim 'revoke-bucket-claim' in namespace 'access-revoke-namespace' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, revokeBucketClaim, true)

		By("Checking if the Bucket referencing 'revoke-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(revokeBucket, true)

		By("Checking if the Bucket referencing 'revoke-bucket-claim' bucketID is not empty")
		steps.CheckBucketID(revokeBucket)

		By("Creating the BucketAccessClass 'revoke-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, revokeBucketAccessClass)

		By("Creating the BucketAccess 'revoke-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, revokeBucketAccess)

		By("Checking if the BucketAccess 'revoke-bucket-access' has status 'accessGranted' set to 'true")
		steps.CheckBucketAccessStatus(ctx, bucketClient, revokeBucketAccess, true)

		By("Creating User '${user}' in account on ObjectScale platform")
		steps.CheckUser(ctx, IAMClient, revokeBucket.Name, Namespace)

		// I need bucket name here that is generated by one of the steps above, I think. Maybe there is a better way to do this.
		resourceARN := objscl.BuildResourceString(ObjectscaleID, ObjectstoreID, revokeBucket.Name)
		principalUsername = objscl.BuildUsername(Namespace, revokeBucket.Name)
		principalARN := objscl.BuildPrincipalString(Namespace, revokeBucket.Name)
		initialPolicy = policy.Document{
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
		finalPolicy = policy.Document{
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

		By("Creating Policy '${policy}' on ObjectScale platform")
		steps.CheckPolicy(ctx, objectscale, initialPolicy, revokeBucket, Namespace)

		By("Checking if BucketAccess resource 'revoke-bucket-access' in namespace 'access-revoke-namespace' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, revokeBucketAccess, principalUsername)

		By("Checking if Secret 'bucket-credentials-1' is created in namespace 'access-revoke-namespace'")
		steps.CheckSecret(ctx, clientset, validSecret)
	})

	It("Successfully revokes access to bucket", func(ctx context.Context) {
		By("Deleting the BucketAccess 'revoke-bucket-access'")
		steps.DeleteBucketAccessResource(ctx, bucketClient, revokeBucketAccess)

		By("Deleting Policy for Bucket referencing BucketClaim 'revoke-bucket-claim' on ObjectScale platform")
		steps.CheckPolicy(ctx, objectscale, finalPolicy, revokeBucket, Namespace)

		By("Deleting User '${user}' in account on ObjectScale platform")
		steps.CheckUserDeleted(ctx, IAMClient, revokeBucket.Name, Namespace)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketClaimResource(ctx, bucketClient, revokeBucketClaim)
			steps.DeleteBucketClassResource(ctx, bucketClient, revokeBucketClass)
		})
	})
})
