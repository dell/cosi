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

package steps

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cosiapi "sigs.k8s.io/container-object-storage-interface-api/apis"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

const (
	// bucketInfo indicates name of data entry in secret, where the all information
	// created by COSI driver is stored.
	bucketInfo = "BucketInfo"

	// testObjectKey is a key of an object that is put and deleted from bucket.
	testObjectKey = "cosi-test.txt"

	// testObjectData is a data of an object that is put and deleted from bucket.
	testObjectData = "COSI test data 💀"
)

// CheckClusterAvailability Ensure that Kubernetes cluster is available.
func CheckClusterAvailability(ctx context.Context, clientset *kubernetes.Clientset) {
	value, err := clientset.ServerVersion()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(value).ToNot(gomega.BeNil())
}

// CreateNamespace Ensure that Kubernetes namespace is created.
func CreateNamespace(ctx context.Context, clientset *kubernetes.Clientset, namespace string) {
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		namespaceObj := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err := clientset.CoreV1().Namespaces().Create(ctx, namespaceObj, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	} else {
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}
}

// DeleteNamespace Ensure that Kubernetes namespace is deleted.
func DeleteNamespace(ctx context.Context, clientset *kubernetes.Clientset, namespace string) {
	err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketClassSpec Ensure that specification of custom resource "my-bucket-class" is correct.
func CheckBucketClassSpec(_ *kubernetes.Clientset, _ v1alpha1.BucketClaimSpec) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckSecret is used to check if secret exists.
func CheckSecret(ctx context.Context, clientset *kubernetes.Clientset, secret *v1.Secret) *v1.Secret {
	sec, err := clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(sec).NotTo(gomega.BeNil())
	gomega.Expect(sec.Name).To(gomega.Equal(secret.Name))
	gomega.Expect(sec.Namespace).To(gomega.Equal(secret.Namespace))
	gomega.Expect(sec.Data).NotTo(gomega.Or(gomega.BeNil(), gomega.BeEmpty()))
	return sec
}

// CheckBucketClaimEvents Check BucketClaim events.
func CheckBucketClaimEvents(ctx context.Context, clientset *kubernetes.Clientset, bucketClaim *v1alpha1.BucketClaim, expected *v1.Event) {
	listOptions := metav1.ListOptions{}

	listOptions.FieldSelector = fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.kind", bucketClaim.Kind),
		fields.OneTermEqualSelector("involvedObject.name", bucketClaim.Name),
		fields.OneTermEqualSelector("type", expected.Type),
	).String()

	eventList := &v1.EventList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EventList",
			APIVersion: "v1",
		},
	}

	for {
		list, err := clientset.CoreV1().Events(bucketClaim.Namespace).List(ctx, listOptions)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		eventList.Items = append(eventList.Items, list.Items...)

		nextContinueToken, _ := meta.NewAccessor().Continue(list)
		if len(nextContinueToken) == 0 {
			break
		}

		listOptions.Continue = nextContinueToken
	}

	// check if there is event having required reason
	gomega.Expect(eventList.Items).NotTo(gomega.BeEmpty())

	found := false

	for _, event := range eventList.Items {
		if event.Reason == expected.Reason {
			found = true
			break
		}
	}

	gomega.Expect(found).To(gomega.BeTrue())
}

func GetAccessKeyID(ctx context.Context, clientset *kubernetes.Clientset, validSecret *v1.Secret) string {
	secret, err := clientset.CoreV1().Secrets(validSecret.Namespace).Get(ctx, validSecret.Name, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	var secretData cosiapi.BucketInfo

	err = json.Unmarshal(secret.Data[bucketInfo], &secretData)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	accessKey := secretData.Spec.S3.AccessKeyID
	return accessKey
}

// CheckBucketAccessFromSecret Check if Bucket can be accessed with data from specified secret.
func CheckBucketAccessFromSecret(ctx context.Context, clientset *kubernetes.Clientset, bucketClient *bucketclientset.Clientset, validSecret *v1.Secret, myBucket *v1alpha1.Bucket) {
	secret, err := clientset.CoreV1().Secrets(validSecret.Namespace).Get(ctx, validSecret.Name, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	var secretData cosiapi.BucketInfo

	err = json.Unmarshal(secret.Data[bucketInfo], &secretData)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	accessKey := secretData.Spec.S3.AccessKeyID
	secretKey := secretData.Spec.S3.AccessSecretKey
	s3Endpoint := secretData.Spec.S3.Endpoint
	bucketName := secretData.Spec.BucketName

	x509Client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(s3Endpoint),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient:       &x509Client,
		Region:           aws.String("us-east-1"),
	}

	session, err := session.NewSession(s3Config)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	s3Client := s3.New(session)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(testObjectData),
		Bucket: aws.String(bucketName),
		Key:    aws.String(testObjectKey),
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(testObjectKey),
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// DeleteSecret first removes the finalizers from secret and then deletes it.
func DeleteSecret(ctx context.Context, clientset *kubernetes.Clientset, secret *v1.Secret) {
	secret, err := clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	secret.SetFinalizers([]string{})

	_, err = clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err = clientset.CoreV1().Secrets(secret.Namespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}
