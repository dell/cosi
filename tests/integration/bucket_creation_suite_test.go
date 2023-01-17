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

func TestBucketCreation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucket Creation Suite")
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
		bucketClaimValid = &v1alpha1.BucketClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bucket-claim-valid",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "my-bucket-class",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
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
		bucketClaimInvalid = &v1alpha1.BucketClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "BucketClaim",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bucket-claim-invalid",
				Namespace: "namespace-1",
			},
			Spec: v1alpha1.BucketClaimSpec{
				BucketClassName: "bucket-class-invalid",
				Protocols: []v1alpha1.Protocol{
					v1alpha1.ProtocolS3,
				},
			},
		}
		validBucket = &v1alpha1.Bucket{
			Spec: v1alpha1.BucketSpec{
				BucketClassName: "my-bucket-class",
				BucketClaim:     &v1.ObjectReference{Kind: "BucketClass", Name: "bucket-claim-valid", Namespace: "namespace-1"},
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
		By("Creating the BucketClass 'my-bucket-class' is created")
		steps.CreateBucketClassResource(bucketClient, myBucketClass)
	})

	// STEP: Scenario: Successfull bucket creation
	It("Successfully creates bucket", func() {
		// STEP: BucketClaim resource is created from specification "bucket-claim-valid"
		By("creating a BucketClaim resource from specification 'bucket-claim-valid'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimValid)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" is created in ObjectStore "object-store-1"
		By("checking if Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, validBucket)
		// STEP: BucketClaim resource "bucket-claim-valid" in namespace "namespace-1" status "bucketReady" is "true"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-valid' in namespace 'namespace-1' is 'true'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimValid)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketReady" is "true" and bucketID is not empty
		By("checking the status 'bucketReady' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid'  is 'true'")
		steps.CheckBucketStatus(bucketClient, validBucket)
		// STEP: Bucket resource referencing BucketClaim resource "bucket-claim-valid" status "bucketID" is not empty
		By("checking the status 'bucketID' of Bucket resource referencing BucketClaim resource 'bucket-claim-valid' is not empty")
		steps.CheckBucketID(bucketClient, validBucket)
	})
	// STEP: Scenario: Unsuccessfull bucket creation
	It("Unsuccessfully tries to create bucket", func() {
		// STEP: BucketClaim resource is created from specification "bucket-claim-invalid"
		By("creating a BucketClaim resource from specification 'bucket-claim-invalid'")
		steps.CreateBucketClaimResource(bucketClient, bucketClaimInvalid)
		// STEP: BucketClaim resource "bucket-claim-invalid" in namespace "namespace-1" status "bucketReady" is "false"
		By("checking if the status 'bucketReady' of BucketClaim resource 'bucket-claim-invalid' in namespace 'namespace-1' is 'false'")
		steps.CheckBucketClaimStatus(bucketClient, bucketClaimInvalid)
	})
})

var _ = AfterSuite(func() {
	// CLean up
})
