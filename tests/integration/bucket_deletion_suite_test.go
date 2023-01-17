package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/dell/cosi-driver/tests/integration/steps"
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBucketDeletion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucket Deletion Suite")
}

var _ = BeforeSuite(func() {
	// Global setup
	// a way to inject k8s conifg from env
	kubeConfig := os.Getenv("KUBECONFIG")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	Expect(err).To(BeNil())

	// k8s clientset
	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).To(BeNil())

	// Bucket clientset
	bucketClient, err = bucketclientset.NewForConfig(cfg)
	Expect(err).To(BeNil())

	// ObjectScale clientset
	// TODO: check how to connect to objectscale with parameters for this function
	objectscale = objectscaleRest.NewClientSet(
		"https://testserver",
		"https://testgateway",
		"svc-objectscale-domain-c8",
		"objectscale-graphql-7d754f8499-ng4h6",
		"OSC234DSF223423",
		"IgQBVjz4mq1M6wmKjHmfDgoNSC56NGPDbLvnkaiuaZKpwHOMFOMGouNld7GXCC690qgw4nRCzj3EkLFgPitA2y8vagG6r3yrUbBdI8FsGRQqW741eiYykf4dTvcwq8P6",
		http.DefaultClient,
		false,
	)
})

var _ = Describe("COSI driver", func() {
	// Resources for scenarios
	var (
		bucketClassDelete = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class-delete",
			},
			DeletionPolicy: "delete",
			DriverName:     "cosi-driver",
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		bucketClassRetain = &v1alpha1.BucketClass{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-bucket-class-retain",
			},
			DeletionPolicy: "retain",
			DriverName:     "cosi-driver",
			Parameters: map[string]string{
				"objectScaleID": "${objectScaleID}",
				"objectStoreID": "${objectStoreID}",
				"accountSecret": "${secretName}",
			},
		}
		bucketClaimDelete = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-claim-delete",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class-delete",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		bucketClaimRetain = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bucket-claim-retain",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class-retain",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		deleteBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class-delete",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "bucket-claim-delete", Namespace: "namespace-1"},
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		retainBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class-retain",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "bucket-claim-retain", Namespace: "namespace-1"},
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
	)

	// Background
	BeforeEach(func() {
		By("Checking if the cluster is ready")
		steps.CheckClusterAvailability(clientset)
		By("Checking if the ObjectScale platform is ready")
		steps.CheckObjectScaleInstallation(clientset)
		By("Checking if the ObjectStore 'object-store-1' is created")
		steps.CreateObjectStore(objectscale, "object-store-1")
		By("Checking if namespace 'driver-ns' is created")
		steps.CreateNamespace(clientset, "driver-ns")
		By("Checking if namespace 'namespace-1' is created")
		steps.CreateNamespace(clientset, "namespace-1")
		By("Checking if COSI controller 'cosi-controller' is installed in namespace 'driver-ns'")
		steps.CheckCOSIControllerInstallation(clientset, "cosi-controller", "driver-ns")
		By("Checking if COSI driver 'cosi-driver' is installed in namespace 'driver-ns'")
		steps.CheckCOSIDriverInstallation(clientset, "cosi-driver", "driver-ns")
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "delete"
	It("Delets the bucket when deletionPolicy is set to 'delete'", func() {
		// STEP: BucketClass resource is created from specification "my-bucket-class-delete"
		By("creating a BucketClass resource from specification 'my-bucket-class-delete'")
		steps.CreateBucketClassResource(bucketClient, bucketClassDelete)

		// STEP: BucketClaim resource is created from specification "my-bucket-claim-delete"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-delete'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimDelete)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" is created in ObjectStore "object-store-1"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, deleteBucket)

		// STEP: BucketClaim resource "bucket-claim-delete" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-delete' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimDelete)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete'  is 'true'")
		steps.CheckBucketStatus(bucketClient, deleteBucket)

		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-delete" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-delete' is not empty")
		steps.CheckBucketStatus(bucketClient, deleteBucket)

		// STEP: BucketClaim resource "my-bucket-claim-delete" is deleted in namespace "namespace-1"
		By("deleting BucketClaim resource 'my-bucket-claim-delete' in namespace 'namespace-1'")
		steps.DeleteBucketClaimResource(bucketClient, bucketClaimDelete)

		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-delete" is deleted in ObjectStore "object-store-1"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-delete' is deleted in ObjectStore 'object-store-1'")
		steps.CheckBucketDeletionInObjectStore(objectscale, deleteBucket)
	})

	// STEP: Scenario: BucketClaim deletion with deletionPolicy set to "retain"
	It("Does not delete the bucket when deletionPolicy is set to 'retain'", func() {
		// STEP: BucketClass resource is created from specification "my-bucket-class-retain"
		By("creating a BucketClass resource from specification 'my-bucket-class-retain'")
		steps.CreateBucketClassResource(bucketClient, bucketClassRetain)
		// STEP: BucketClaim resource is created from specification "my-bucket-claim-retain"
		By("creating a BucketClaim resource from specification 'my-bucket-claim-retain'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimRetain)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" is created in ObjectStore "object-store-1"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, retainBucket)
		// STEP: BucketClaim resource "bucket-claim-retain" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-retain' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimRetain)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-retain'  is 'true'")
		steps.CheckBucketStatus(bucketClient, retainBucket)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-retain" status "bucketID" is not empty
		By("checking the ID of Bucket resource referencing BucketClaim resource 'bucket-claim-retain' is not empty")
		// STEP: BucketClaim resource "my-bucket-claim-retain" is deleted in namespace "namespace-1"
		By("deleting BucketClaim resource 'my-bucket-claim-retain' in namespace 'namespace-1'")
		steps.DeleteBucketClaimResource(bucketClient, bucketClaimRetain)
		// STEP: Bucket referencing BucketClaim resource "my-bucket-claim-retain" is available in ObjectStore "object-store-1"
		By("checking if Bucket referencing BucketClaim resource 'my-bucket-claim-retain' is available in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, retainBucket)
	})

})

var _ = AfterSuite(func() {
	// CLean up
})
