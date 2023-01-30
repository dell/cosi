//go:build integration

package main_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	objectscaleClient "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// place for storing global variables like specs
var (
	clientset    *kubernetes.Clientset
	bucketClient *bucketclientset.Clientset
	objectscale  *objectscaleRest.ClientSet
	iamClient    *iam.IAM
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "COSI Integration Suite")
}

var _ = BeforeSuite(func() {
	// Global setup
	// Load environment variables

	kubeConfig, exists := os.LookupEnv("KUBECONFIG")
	Expect(exists).To(BeTrue())
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	Expect(err).To(BeNil())

	objectscaleURL, exists := os.LookupEnv("OBJECTSCALE_URL")
	Expect(exists).To(BeTrue())

	objectscaleUser, exists := os.LookupEnv("OBJECTSCALE_USER")
	Expect(exists).To(BeTrue())

	objectscalePassword, exists := os.LookupEnv("OBJECTSCALE_PASSWORD")
	Expect(exists).To(BeTrue())

	// k8s clientset
	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).To(BeNil())

	// Bucket clientset
	bucketClient, err = bucketclientset.NewForConfig(cfg)
	Expect(err).To(BeNil())

	// ObjectScale clientset
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	unsafeClient := &http.Client{Transport: transport}
	objectscaleLoginAddress := fmt.Sprintf("%s:31613", objectscaleURL)
	objectscaleManagementAddress := fmt.Sprintf("%s:30007", objectscaleURL)

	objectscale = objectscaleRest.NewClientSet(
		objectscaleClient.NewClient(
			objectscaleManagementAddress,
			objectscaleLoginAddress,
			objectscaleUser,
			objectscalePassword,
			unsafeClient,
			false,
		),
	)

	// IAM clientset
	var (
		endpoint = objectscaleLoginAddress
		region   = "us-west-2"
	)
	iamSession, err := session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	})
	Expect(err).To(BeNil())
	iamClient = iam.New(iamSession)
})

var _ = AfterSuite(func() {
	// Global teardown
})
