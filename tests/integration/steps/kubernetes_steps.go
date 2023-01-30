package steps

import (
	"context"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// CheckClusterAvailability Ensure that Kubernetes cluster is available
func CheckClusterAvailability(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset) {
	value, err := clientset.ServerVersion()
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(value).ToNot(gomega.BeNil())
}

// CreateNamespace Ensure that Kubernetes namespace is created
func CreateNamespace(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, namespace string) {
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		namespaceObj := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err := clientset.CoreV1().Namespaces().Create(ctx, namespaceObj, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		gomega.Expect(err).To(gomega.BeNil())
	}
}

// DeleteNamespace Ensure that Kubernetes namespace is deleted
func DeleteNamespace(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, namespace string) {
	err := clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	gomega.Expect(err).To(gomega.BeNil())
}

// CheckBucketClassSpec Ensure that specification of custom resource "my-bucket-class" is correct
func CheckBucketClassSpec(clientset *kubernetes.Clientset, bucketClassSpec v1alpha1.BucketClaimSpec) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// Check if secret exists
func CheckSecret(ctx context.Context, clientset *kubernetes.Clientset, secret *v1.Secret) {
	sec, err := clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	if gomega.Expect(sec).NotTo(gomega.BeNil()) {
		gomega.Expect(sec.Name).To(gomega.Equal(secret.Namespace))
		gomega.Expect(sec.Namespace).To(gomega.Equal(secret.Namespace))
		gomega.Expect(sec.Data).NotTo(gomega.Or(gomega.BeNil(), gomega.BeEmpty()))
	}
}

// CheckBucketClaimEvents Check BucketClaim events
func CheckBucketClaimEvents(ctx context.Context, clientset *kubernetes.Clientset, bucketClaim *v1alpha1.BucketClaim, expected string) {
	el, err := clientset.EventsV1().Events(bucketClaim.Namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + bucketClaim.Name, // FIXME: this is not valid, and fails
	})
	gomega.Expect(err).To(gomega.BeNil())
	if gomega.Expect(el).NotTo(gomega.Or(gomega.BeNil(), gomega.BeEmpty())) {
		found := false
		for _, event := range el.Items {
			if event.Reason == expected {
				found = true
				break
			}
		}
		gomega.Expect(found).To(gomega.Equal(true))
	}
}
