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
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	objectscaleClient "github.com/dell/goobjectscale/pkg/client/rest/client"
	objectscaleIAM "github.com/dell/goobjectscale/pkg/client/rest/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// place for storing global variables like specs.
var (
	clientset     *kubernetes.Clientset
	bucketClient  *bucketclientset.Clientset
	objectscale   *objectscaleRest.ClientSet
	IAMClient     *iam.IAM
	Namespace     string
	ObjectstoreID string
	ObjectscaleID string
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
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
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
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

var _ = AfterSuite(func() {
	// Global teardown
})
