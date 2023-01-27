// go:build integration

package main_test

import (
	"github.com/dell/cosi-driver/tests/integration/steps"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Bucket Access Revoke", Label("revoke"), func() {
	// Resources for scenarios
	var (
		myBucketClass       *v1alpha1.BucketClass
		myBucketClaim       *v1alpha1.BucketClaim
		myBucket            *v1alpha1.Bucket
		myBucketAccessClass *v1alpha1.BucketAccessClass
		myBucketAccess      *v1alpha1.BucketAccess
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
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		myBucketClaim = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
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
		myBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "my-bucket-claim", Namespace: "namespace-1"},
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
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
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

		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(ctx, clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(ctx, clientset)

		// STEP: ObjectStore "object-store-1" is created
		By("Checking if the ObjectStore 'object-store-1' is created")
		steps.CreateObjectStore(ctx, objectscale, "object-store-1")

		// STEP: Kubernetes namespace "driver-ns" is created
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(ctx, clientset, "driver-ns")

		// STEP: Kubernetes namespace "namespace-1" is created
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(ctx, clientset, "namespace-1")

		// STEP: COSI controller "cosi-controller" is installed in namespace "driver-ns"
		By("Checking if COSI controller 'cosi-controller' is installed in namespace 'driver-ns'")
		steps.CheckCOSIControllerInstallation(clientset, "cosi-controller", "driver-ns")

		// STEP: COSI driver "cosi-driver" is installed in namespace "driver-ns"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(clientset, "cosi-driver", "driver-ns")

		// STEP: BucketClass resource is created from specification "my-bucket-class"
		By("Creating the BucketClass 'my-bucket-class'")
		steps.CreateBucketClassResource(bucketClient, myBucketClass)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim"
		By("Creating the BucketClaim 'my-bucket-claim'")
		steps.CreateBucketClaimResource(bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "object-store-1"
		By("Checking if the Bucket referencing 'my-bucket-claim' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, myBucket)

		// STEP: BucketClaim resource "my-bucket-claim" in namespace "namespace-1" status "bucketReady" is "true"
		By("Checking if the BucketClaim 'my-bucket-claim' in namespace 'namespace-1' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
		By("Checking if the Bucket referencing 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(bucketClient, myBucket)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" bucketID is not empty
		By("Checking if the Bucket referencing 'my-bucket-claim' bucketID is not empty")
		steps.CheckBucketID(bucketClient, myBucket)

		// STEP: BucketAccessClass resource is created from specification "my-bucket-access-class"
		By("Creating the BucketAccessClass 'my-bucket-access-class'")
		steps.CreateBucketAccessClassResource(bucketClient, myBucketAccessClass)

		// STEP: BucketAccess resource is created from specification "my-bucket-access"
		By("Creating the BucketAccess 'my-bucket-access'")
		steps.CreateBucketAccessResource(bucketClient, myBucketAccess)

		// STEP: BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accessGranted" is "true"
		By("Checking if the BucketAccess 'my-bucket-access' has status 'accessGranted' set to 'true")
		steps.CheckBucketAccessStatus(bucketClient, myBucketAccess)

		// STEP: User "${user}" in account on ObjectScale platform is created
		By("Creating User '${user}' in account on ObjectScale platform")
		steps.CreateUser(objectscale, "${user}")

		// STEP: Policy "${policy}" on ObjectScale platform is created
		By("Creating Policy '${policy}' on ObjectScale platform")
		steps.CreatePolicy(objectscale, "${policy}", myBucket)

		// STEP: BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accountID" is "${accountID}"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(bucketClient, myBucketAccess, "${accountID}")

		// STEP: Secret "bucket-credentials-1" is created in namespace "namespace-1" and is not empty
		By("Checking if Secret ''bucket-credentials-1' is created in namespace 'namespace-1'")
		steps.CheckSecret(clientset, "bucket-credentials-1", "namespace-1")

		DeferCleanup(func() {
			// Cleanup for background
		})
	})

	// STEP: Revoke access to bucket
	It("Successfully revokes access to bucket", func() {
		// STEP: BucketAccess resource "my-bucket-access" in namespace "namespace-1" is deleted
		By("Deleting the BucketAccess 'my-bucket-access'")
		steps.DeleteBucketAccessResource(bucketClient, myBucketAccess)

		// STEP: Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is deleted
		By("Deleting Policy for Bucket referencing BucketClaim 'my-bucket-claim' on ObjectScale platform")
		steps.DeletePolicy(objectscale, myBucket)

		// STEP: User "${user}" in account on ObjectScale platform is deleted
		By("Deleting User '${user}' in account on ObjectScale platform")
		steps.DeleteUser(objectscale, "${user}")

		DeferCleanup(func() {
			// Cleanup for scenario: Revoke access to bucket
		})
	})
})
