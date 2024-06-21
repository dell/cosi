// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build integration

package main_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dell/cosi/tests/integration/steps"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	objectscaleClient "github.com/dell/goobjectscale/pkg/client/rest/client"
	objectscaleIAM "github.com/dell/goobjectscale/pkg/client/rest/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// place for storing global variables like specs.
var (
	clientset    *kubernetes.Clientset
	bucketClient *bucketclientset.Clientset
	objectscale  *objectscaleRest.ClientSet
	IAMClient    *iam.IAM

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

	objectstoreGateway, exists := os.LookupEnv("OBJECTSCALE_OBJECTSTORE_GATEWAY")
	Expect(exists).To(BeTrue())

	objectscaleUser, exists := os.LookupEnv("OBJECTSCALE_USER")
	Expect(exists).To(BeTrue())

	objectscalePassword, exists := os.LookupEnv("OBJECTSCALE_PASSWORD")
	Expect(exists).To(BeTrue())

	ObjectstoreID, exists = os.LookupEnv("OBJECTSCALE_OBJECTSTORE_ID")
	Expect(exists).To(BeTrue())

	ObjectscaleID, exists = os.LookupEnv("OBJECTSCALE_ID")
	Expect(exists).To(BeTrue())

	// k8s clientset
	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())

	// Bucket clientset
	bucketClient, err = bucketclientset.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())

	// ObjectScale clientset
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			CipherSuites:       getSecuredCipherSuites(),
		},
	}
	unsafeClient := &http.Client{Transport: transport}

	objectscaleAuthUser := objectscaleClient.AuthUser{
		Gateway:  objectscaleGateway,
		Username: objectscaleUser,
		Password: objectscalePassword,
	}
	objectscale = objectscaleRest.NewClientSet(
		&objectscaleClient.Simple{
			Endpoint:       objectstoreGateway,
			Authenticator:  &objectscaleAuthUser,
			HTTPClient:     unsafeClient,
			OverrideHeader: false,
		},
	)

	// IAM clientset
	var (
		endpoint = objectscaleGateway + "/iam"
		region   = "us-west-2"
	)
	iamSession, err := session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint:gosec
					CipherSuites:       getSecuredCipherSuites(),
				},
			},
		},
	})

	Expect(err).ToNot(HaveOccurred())

	IAMClient = iam.New(iamSession)
	err = objectscaleIAM.InjectTokenToIAMClient(IAMClient, &objectscaleAuthUser, *unsafeClient)
	Expect(err).ToNot(HaveOccurred())
	err = objectscaleIAM.InjectAccountIDToIAMClient(IAMClient, Namespace)
	Expect(err).ToNot(HaveOccurred())
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
