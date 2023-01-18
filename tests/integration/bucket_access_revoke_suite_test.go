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

func TestBucketAccessRevoke(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucket Access Revoke Suite")
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

		By("Creating the BucketClass 'my-bucket-class'")
		steps.CreateBucketClassResource(bucketClient, myBucketClass)

		By("Creating the BucketClaim 'my-bucket-claim'")
		steps.CreateBucketClaimResource(bucketClient, myBucketClaim)

		By("Checking if the Bucket referencing 'my-bucket-claim' is created in ObjectStore 'object-store-1'")
		steps.CheckBucketResourceInObjectStore(objectscale, myBucket)

		By("Checking if the Bucket referencing 'my-bucket-claim' is created in namespace 'namespace-1'")
		steps.CheckBucketClaimStatus(bucketClient, myBucketClaim)

		By("Checking if the Bucket referencing 'my-bucket-claim' bucketID is not empty")
		steps.CheckBucketID(bucketClient, myBucket)

		By("Creating the BucketAccessClass 'my-bucket-access-class'")
		steps.CreateBucketAccessClassResource(bucketClient, myBucketAccessClass)

		By("Creating the BucketAccess 'my-bucket-access'")
		steps.CreateBucketAccessResource(bucketClient, myBucketAccess)

		By("Checking if the BucketAccess 'my-bucket-access' has status 'accessGranted' set to 'true")
		steps.CheckBucketAccessStatus(bucketClient, myBucketAccess)

		By("Creating User '${user}' in account on ObjectScale platform")
		steps.CreateUser(objectscale, "${user}")

		By("Creating Policy '${policy}' on ObjectScale platform")
		steps.CreatePolicy(objectscale, "${policy}", myBucket)

		By("Checking if BucketAccess resource 'my-bucket-access' in namespace 'namespace-1' status 'accountID' is '${accountID}'")
		steps.CheckBucketAccessAccountID(bucketClient, myBucketAccess, "${accountID}")

		By("Checking if Secret ''bucket-credentials-1' is created in namespace 'namespace-1'")
		steps.CheckSecret(clientset, "bucket-credentials-1", "namespace-1")
	})

	// STEP: Revoke access to bucket
	It("Successfully revokes access to bucket", func() {
		By("Deleting the BucketAccess 'my-bucket-access'")
		steps.DeleteBucketAccessResource(bucketClient, myBucketAccess)

		By("Deleting Policy for Bucket referencing BucketClaim 'my-bucket-claim' on ObjectScale platform")
		steps.DeletePolicy(objectscale, myBucket)

		By("Deleting User '${user}' in account on ObjectScale platform")
		steps.DeleteUser(objectscale, "${user}")
	})
})

var _ = AfterSuite(func() {
	// CLean up
})
