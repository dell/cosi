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
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// CheckClusterAvailability Ensure that Kubernetes cluster is available.
func CheckClusterAvailability(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset) {
	value, err := clientset.ServerVersion()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(value).ToNot(gomega.BeNil())
}

// CreateNamespace Ensure that Kubernetes namespace is created.
func CreateNamespace(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, namespace string) {
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
func DeleteNamespace(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, namespace string) {
	err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketClassSpec Ensure that specification of custom resource "my-bucket-class" is correct.
func CheckBucketClassSpec(clientset *kubernetes.Clientset, bucketClassSpec v1alpha1.BucketClaimSpec) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckSecret is used to check if secret exists.
func CheckSecret(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, secret *v1.Secret) {
	sec, err := clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(sec).NotTo(gomega.BeNil())
	gomega.Expect(sec.Name).To(gomega.Equal(secret.Namespace))
	gomega.Expect(sec.Namespace).To(gomega.Equal(secret.Namespace))
	gomega.Expect(sec.Data).NotTo(gomega.Or(gomega.BeNil(), gomega.BeEmpty()))
}

// CheckBucketClaimEvents Check BucketClaim events.
func CheckBucketClaimEvents(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, bucketClaim *v1alpha1.BucketClaim, expected *v1.Event) {
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
