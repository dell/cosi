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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	smithyhttp "github.com/aws/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/smithy-go/middleware"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dell/cosi/tests/integration/steps"
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/rest"
	"github.com/dell/goobjectscale/pkg/client/rest/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	bucketclientset "sigs.k8s.io/container-object-storage-interface/client/clientset/versioned"
)

// place for storing global variables like specs.
var (
	clientset    *kubernetes.Clientset
	bucketClient *bucketclientset.Clientset

	mgmtClient api.ClientSet
	IAMClient  *iam.Client

	Namespace     string
	ObjectstoreID string
	ObjectscaleID string

	DeploymentName      string
	DriverNamespace     string
	DriverContainerName string
)

const (
	DriverID = "e2e.test.objectscale"
)

func TestIntegration(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "COSI Integration Suite")
}

var _ = BeforeSuite(func() {
	// Global setup
	// Load environment variables

	exists := false

	DeploymentName, exists = os.LookupEnv("HELM_RELEASE_NAME")
	Expect(exists).To(BeTrue())

	DriverNamespace, exists = os.LookupEnv("DRIVER_NAMESPACE")
	Expect(exists).To(BeTrue())

	DriverContainerName, exists = os.LookupEnv("DRIVER_CONTAINER_NAME")
	Expect(exists).To(BeTrue())

	Namespace, exists = os.LookupEnv("OBJECTSCALE_NAMESPACE")
	Expect(exists).To(BeTrue())

	kubeConfig, exists := os.LookupEnv("KUBECONFIG")
	Expect(exists).To(BeTrue())
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	Expect(err).ToNot(HaveOccurred())

	objectscaleGateway, exists := os.LookupEnv("OBJECTSCALE_GATEWAY")
	Expect(exists).To(BeTrue())

	objectscaleUser, exists := os.LookupEnv("OBJECTSCALE_USER")
	Expect(exists).To(BeTrue())

	objectscalePassword, exists := os.LookupEnv("OBJECTSCALE_PASSWORD")
	Expect(exists).To(BeTrue())

	// k8s clientset
	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())

	// Bucket clientset
	bucketClient, err = bucketclientset.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())

	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	httpClient := &http.Client{
		Transport: baseTransport,
	}

	objectscaleAuthUser := client.AuthUser{
		Gateway:  objectscaleGateway,
		Username: objectscaleUser,
		Password: objectscalePassword,
	}

	mgmtClient = rest.NewClientSet(&client.Simple{
		Endpoint:       objectscaleGateway,
		Authenticator:  &objectscaleAuthUser,
		OverrideHeader: false,
		HTTPClient:     httpClient,
	})

	Expect(err).ToNot(HaveOccurred())

	// login so we can use token for IAM
	err = objectscaleAuthUser.Login(context.Background(), httpClient)
	Expect(err).ToNot(HaveOccurred())

	config, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"root", "---", "")),
		awsconfig.WithHTTPClient(httpClient),
		awsconfig.WithAPIOptions([]func(*middleware.Stack) error{
			smithyhttp.AddHeaderValue("X-Emc-Namespace", Namespace),
			smithyhttp.AddHeaderValue("X-Sds-Auth-Token", objectscaleAuthUser.Token()),
		}),
	)
	Expect(err).ToNot(HaveOccurred())

	IAMClient = iam.NewFromConfig(config, func(o *iam.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("%s/iam/", objectscaleGateway))
	})
})

var _ = AfterSuite(func(ctx context.Context) {
	podList, err := clientset.CoreV1().Pods(DriverNamespace).List(ctx, v1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(podList.Items).ToNot(BeEmpty())

	// ensure every job from test is processed
	time.Sleep(time.Second)
	for _, pod := range podList.Items {
		steps.CheckErrors(ctx, clientset, pod.Name, DriverContainerName, pod.Namespace)
	}
})

func getSecuredCipherSuites() (suites []uint16) {
	securedSuite := tls.CipherSuites()
	for _, v := range securedSuite {
		suites = append(suites, v.ID)
	}

	return suites
}
