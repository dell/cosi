package main_test

import (
	"github.com/dell/cosi-driver/tests/integration/steps"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Bucket Access KEY", Label("key-flow"), func() {
	// Resources for scenarios
	var (
		myBucketClass = &v1alpha1.BucketClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class",
			},
			DriverName:     "cosi-driver",
			DeletionPolicy: v1alpha1.DeletionPolicyDelete,
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
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
		myBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClaim", Name: "my-bucket-claim", Namespace: "namespace-1"},
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
	)

	// Background
	BeforeEach(func() {
		// STEP: Kubernetes cluster is up and running
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)

		// STEP: ObjectScale platform is installed on the cluster
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(clientset)

		// STEP: ObjectStore "object-store-1" is created
		By("Checking if the ObjectStore 'object-store-1' is created")
		steps.CreateObjectStore(objectscale, "object-store-1")

		// STEP: Kubernetes namespace "driver-ns" is created
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(clientset, "driver-ns")

		// STEP: Kubernetes namespace "namespace-1" is created
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(clientset, "namespace-1")

		// STEP: COSI controller "cosi-controller" is installed in namespace "driver-ns"
		By("Checking if COSI controller 'cosi-controller' is installed in namespace 'driver-ns'")
		steps.CheckCOSIControllerInstallation(clientset, "cosi-controller", "driver-ns")

		// STEP: COSI driver "cosi-driver" is installed in namespace "driver-ns"
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(clientset, "cosi-driver", "driver-ns")

		// STEP: BucketClass resource is created from specification "my-bucket-class"
		By("Creating the BucketClass 'my-bucket-class' is created")
		steps.CreateBucketClassResource(bucketClient, myBucketClass)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim"
		By("Creating the BucketAccessClass 'my-bucket-access-class' is created")
		steps.CreateBucketClaimResource(bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" is created in ObjectStore "object-store-1"
		By("Checking if bucket referencing 'my-bucket-claim' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, myBucket)

		// STEP: BucketClaim resource "my-bucket-claim" in namespace "namespace-1" status "bucketReady" is "true"
		By("Checking if BucketClaim resource 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, myBucketClaim)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket-claim" status "bucketReady" is "true"
		By("Checking if Bucket resource referencing 'my-bucket-claim' status 'bucketReady' is 'true'")
		steps.CheckBucketStatus(bucketClient, myBucket)

		// STEP: Bucket resource referencing BucketClaim resource "my-bucket" bucketID is not empty
		By("Checking if Bucket resource 'my-bucket' status 'bucketID' is not empty")
		steps.CheckBucketID(bucketClient, myBucket)
	})

	// STEP: Scenario: BucketAccess creation with KEY authorization mechanism
	It("Creates BucketAccess with KEY authorization mechanism", func() {
		// STEP: BucketAccessClass resource is created from specification "my-bucket-access-class"
		By("Creating BucketAccessClass resource 'my-bucket-access-class'")
		steps.CreateBucketAccessClassResource(bucketClient, myBucketAccessClass)

		// STEP: BucketAccess resource is created from specification "my-bucket-access"
		By("Creating BucketAccess resource 'my-bucket-access'")
		steps.CreateBucketAccessResource(bucketClient, myBucketAccess)

		// STEP: BucketAccess resource "my-bucket-access" status "accessGranted" is "true"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accessGranted' is 'true'")
		steps.CheckBucketAccessStatus(bucketClient, myBucketAccess)

		// STEP: User "${user}" in account on ObjectScale platform is created
		By("Creating User resource '${user}'")
		steps.CreateUser(objectscale, "${user}")

		// STEP: Policy "${policy}" for Bucket resource referencing BucketClaim resource "my-bucket-claim" on ObjectScale platform is created
		By("Creating Policy resource '${policy}' for Bucket resource referencing BucketClaim resource 'my-bucket-claim'")
		steps.CreatePolicy(objectscale, "${policy}", myBucket)

		// STEP: BucketAccess resource "my-bucket-access" in namespace "namespace-1" status "accountID" is "${accountID}"
		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(bucketClient, myBucketAccess, "${accountID}")

		// STEP: Secret "bucket-credentials-1" is created in namespace "namespace-1" and is not empty
		By("Checking if Secret 'bucket-credentials-1' in namespace 'namespace-1' is not empty")
		steps.CheckSecret(clientset, "bucket-credentials-1", "namespace-1")

		//STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is accessible from Secret "bucket-credentials-1"
		By("Checking if Bucket resource referencing BucketClaim resource 'my-bucket-claim' is accessible from Secret 'bucket-credentials-1'")
		steps.CheckBucketAccessFromSecret(objectscale, myBucket, "bucket-credentials-1")
	})
})
