// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"

	. "github.com/onsi/ginkgo/v2"

	v1 "k8s.io/api/core/v1"

	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/cosi/tests/integration/steps"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Bucket Access Grant for Multiple Greenfield Buckets", Ordered, Label("grant", "greenfield", "objectscale", "multiple"), func() {
	// Resources for scenarios
	const (
		namespace string = "access-grant-multiple-namespace-greenfield"
	)

	var (
		grantBucketClass       *v1alpha1.BucketClass
		grantBucketClaim       *v1alpha1.BucketClaim
		grantBucket            *v1alpha1.Bucket
		grantBucketAccessClass *v1alpha1.BucketAccessClass
		grantBucketAccess1     *v1alpha1.BucketAccess
		grantBucketAccess2     *v1alpha1.BucketAccess
		validSecret1           *v1.Secret
		validSecret2           *v1.Secret
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
				Name: "greenfield-grant-multiple-bucket-class",
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
				Name:      "greenfield-grant-multiple-bucket-claim",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "greenfield-grant-multiple-bucket-class",
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
				Name: "greenfield-grant-multiple-bucket-access-class",
			},
			DriverName:         "cosi.dellemc.com",
			AuthenticationType: v1alpha1.AuthenticationTypeKey,
			Parameters: map[string]string{
				"id": DriverID,
			},
		}
		grantBucketAccess1 = &v1alpha1.BucketAccess{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccess",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "greenfield-grant-multiple-bucket-access-1",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "greenfield-grant-multiple-bucket-access-class",
				BucketClaimName:       "greenfield-grant-multiple-bucket-claim",
				CredentialsSecretName: "greenfield-grant-multiple-bucket-credentials-1",
			},
		}

		grantBucketAccess2 = &v1alpha1.BucketAccess{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketAccess",
				APIVersion: "objectstorage.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "greenfield-grant-multiple-bucket-access-2",
				Namespace: namespace,
			},
			Spec: v1alpha1.BucketAccessSpec{
				BucketAccessClassName: "greenfield-grant-multiple-bucket-access-class",
				BucketClaimName:       "greenfield-grant-multiple-bucket-claim",
				CredentialsSecretName: "greenfield-grant-multiple-bucket-credentials-2",
			},
		}
		validSecret1 = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "greenfield-grant-multiple-bucket-credentials-1",
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"BucketInfo": []byte(`{"metadata":{"name":""},"spec":{"bucketName":"","authenticationType":"","secretS3":{"endpoint":"","region":"","accessKeyID":"","accessSecretKey":""},"protocols":[]}}`),
			},
		}
		validSecret2 = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "greenfield-grant-multiple-bucket-credentials-2",
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

		By("Checking if namespace 'access-grant-namespace' is created")
		steps.CreateNamespace(ctx, clientset, namespace)

		By("Checking if COSI controller 'objectstorage-controller' is installed in namespace 'default'")
		steps.CheckCOSIControllerInstallation(ctx, clientset, "container-object-storage-controller", "container-object-storage-system")

		By("Checking if COSI driver 'cosi' is installed in namespace 'cosi-test-ns'")
		steps.CheckCOSIDriverInstallation(ctx, clientset, DeploymentName, DriverNamespace)

		By("Creating the BucketClass 'grant-bucket-class' is created")
		grantBucketClass = steps.CreateBucketClassResource(ctx, bucketClient, grantBucketClass)

		By("Creating the BucketClaim 'grant-bucket-claim'")
		steps.CreateBucketClaimResource(ctx, bucketClient, grantBucketClaim)

		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-access-class' is created")
		grantBucket = steps.GetBucketResource(ctx, bucketClient, grantBucketClaim)

		By("Checking if bucket referencing 'grant-bucket-claim' is created in ObjectStore '${objectstoreName}'")
		steps.CheckBucketResourceInObjectStore(ctx, mgmtClient, Namespace, grantBucket)

		By("Checking if BucketClaim resource 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(ctx, bucketClient, grantBucketClaim, true)

		By("Checking if Bucket resource referencing 'grant-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(grantBucket, true)

		By("Checking if Bucket resource 'grant-bucket' status 'bucketID' is not empty")
		steps.CheckBucketID(grantBucket)
	})

	It("Creates BucketAccess with KEY authorization mechanism", func(ctx context.Context) {
		By("Creating BucketAccessClass resource 'grant-bucket-access-class'")
		steps.CreateBucketAccessClassResource(ctx, bucketClient, grantBucketAccessClass)

		// grant access to first bucket
		By("Creating BucketAccess resource 'grant-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, grantBucketAccess1)

		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accessGranted' is 'true'")
		grantBucketAccess1 = steps.CheckBucketAccessStatus(ctx, bucketClient, grantBucketAccess1, true)

		By("Checking if IAM user is created on ObjectScale platform")
		steps.CheckUser(ctx, IAMClient, grantBucketAccess1.Status.AccountID)

		// grant access to second bucket
		By("Creating BucketAccess resource 'grant-bucket-access'")
		steps.CreateBucketAccessResource(ctx, bucketClient, grantBucketAccess2)

		By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accessGranted' is 'true'")
		grantBucketAccess2 = steps.CheckBucketAccessStatus(ctx, bucketClient, grantBucketAccess2, true)

		By("Checking if User ObjectScale platform is created")
		steps.CheckUser(ctx, IAMClient, grantBucketAccess2.Status.AccountID)

		resourceARNBucket := objscl.BuildResourceStrings(grantBucket.Name)
		awsPrincipalStringBucket1 := objscl.BuildPrincipalString(grantBucketAccess1.Status.AccountID, Namespace)
		awsPrincipalStringBucket2 := objscl.BuildPrincipalString(grantBucketAccess2.Status.AccountID, Namespace)

		grantBucketPolicy := GetBucketPolicy(resourceARNBucket,
			[]string{awsPrincipalStringBucket1, awsPrincipalStringBucket2})

		// check bucket policy
		By("Checking if Policy for Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is created")
		steps.CheckPolicy(ctx, mgmtClient, grantBucketPolicy, grantBucket, Namespace)

		// By("Checking if BucketAccess resource 'grant-bucket-access' in namespace 'access-grant-namespace' status 'accountID' is '${accountID}'")
		// steps.CheckBucketAccessAccountID(ctx, bucketClient, grantBucketAccess1, principalUsername1)

		By("Checking if Secret 'bucket-credentials-1' in namespace 'access-grant-namespace' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret1)

		By("Checking if Secret 'bucket-credentials-2' in namespace 'access-grant-namespace' is not empty")
		steps.CheckSecret(ctx, clientset, validSecret2)

		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(ctx, clientset, validSecret1)

		By("Checking if Bucket resource referencing BucketClaim resource 'grant-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(ctx, clientset, validSecret2)

		By("Delete only first bucket access")
		steps.DeleteBucketAccessResource(ctx, bucketClient, grantBucketAccess1)

		grantBucketPolicy = GetBucketPolicy(resourceARNBucket,
			[]string{awsPrincipalStringBucket2})

		By("Checking only second bucket access exists in the policy")
		steps.CheckPolicy(ctx, mgmtClient, grantBucketPolicy, grantBucket, Namespace)

		DeferCleanup(func(ctx context.Context) {
			steps.DeleteBucketAccessResource(ctx, bucketClient, grantBucketAccess2)
			steps.DeleteBucketAccessClassResource(ctx, bucketClient, grantBucketAccessClass)

			// We should wait until access to policy is revoked, to not cause errors.
			steps.CheckEmptyPolicy(ctx, mgmtClient, grantBucket, Namespace)

			steps.DeleteBucketClaimResource(ctx, bucketClient, grantBucketClaim)
			steps.DeleteBucketClassResource(ctx, bucketClient, grantBucketClass)
		})
	})
})

func GetBucketPolicy(resourceARN []string, awsPrincipal []string) policy.Document {
	bucketPolicy := policy.Document{
		Version:   "2012-10-17",
		ID:        "bucket-policy",
		Statement: []policy.StatementEntry{},
	}

	for i := range awsPrincipal {
		bucketPolicy.Statement = append(bucketPolicy.Statement, policy.StatementEntry{
			Effect: "Allow",
			Action: []string{
				"*",
			},
			Resource: resourceARN,
			Principal: map[string]string{
				"AWS": awsPrincipal[i],
			},
			Sid: objscl.PolicySid,
		})
	}
	return bucketPolicy
}
