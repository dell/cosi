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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"
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
		validSecret             *v1.Secret
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
		steps.CheckObjectScaleInstallation(ctx, mgmtClient, Namespace)

		By("Checking if namespace 'cosi-test-ns' is created")
		steps.CreateNamespace(ctx, clientset, DriverNamespace)

		By("Checking if namespace 'access-revoke-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "container-object-storage-controller", "container-object-storage-system")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)

		By("Creating the BucketClass 'revoke-bucket-class'")
		steps.CreateBucketClassResource(ctx, bucketClient, revokeBucketClass)

		By("Creating the BucketClaim 'revoke-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, revokeBucketClaim)

		By("Checking if Bucket resource referencing BucketClaim resource 'revoke-bucket-claim' is created")
		revokeBucket = steps.GetBucketResource(ctx, bucketClient, revokeBucketClaim)

		By("Checking if the Bucket referencing 'revoke-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, revokeBucket)

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
		revokeBucketAccess = steps.CheckBucketAccessStatus(ctx, bucketClient, revokeBucketAccess, true)

		By("Creating User '${user}' in account on ObjectScale platform")
		steps.CheckUser(ctx, IAMClient, revokeBucketAccess.Status.AccountID)

		resourceARNs := objscl.BuildResourceStrings(revokeBucket.Name)
		principalARN := objscl.BuildPrincipalString(revokeBucketAccess.Status.AccountID, Namespace)

		initialPolicy := GetBucketPolicy(resourceARNs, []string{principalARN})

		By("Creating Policy '${policy}' on ObjectScale platform")
		steps.CheckPolicy(ctx, mgmtClient, initialPolicy, revokeBucket, Namespace)

		By("Checking if BucketAccess resource 'revoke-bucket-access' in namespace 'access-revoke-namespace' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(ctx, bucketClient, revokeBucketAccess, revokeBucketAccess.Status.AccountID)

		By("Checking if Secret 'bucket-credentials-1' is created in namespace 'access-revoke-namespace'")
		steps.CheckSecret(ctx, clientset, validSecret)
	})

	It("Successfully revokes access to bucket", func(ctx context.Context) {
		By("Deleting the BucketAccess 'revoke-bucket-access'")
		steps.DeleteBucketAccessResource(ctx, bucketClient, revokeBucketAccess)

		steps.CheckEmptyPolicy(ctx, mgmtClient, revokeBucket, Namespace)

		By("Deleting User '${user}' in account on ObjectScale platform")
		steps.CheckUserDeleted(ctx, IAMClient, revokeBucketAccess.Status.AccountID, Namespace)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketAccessClassResource(ctx, bucketClient, revokeBucketAccessClass)
			steps.DeleteBucketClaimResource(ctx, bucketClient, revokeBucketClaim)
			steps.DeleteBucketClassResource(ctx, bucketClient, revokeBucketClass)
		})
	})
})
