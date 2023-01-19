package main_test

import (
	"net/http"
	"os"
	"testing"

	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
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
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "COSI Integration Suite")
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

var _ = AfterSuite(func() {
	// Global teardown
})
