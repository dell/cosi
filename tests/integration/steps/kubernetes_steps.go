package steps

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// CheckClusterAvailability Ensure that Kubernetes cluster is available
func CheckClusterAvailability(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset) {
	value, err := clientset.ServerVersion()
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(value).ShouldNot(gomega.BeNil())
}

// CreateNamespace Ensure that Kubernetes namespace is created
func CreateNamespace(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, namespace string) {
	namespaceObj := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	_, err := clientset.CoreV1().Namespaces().Create(ctx, namespaceObj, metav1.CreateOptions{})
	gomega.Expect(err).Should(gomega.BeNil())
}

// CheckBucketClassSpec Ensure that specification of custom resource "my-bucket-class" is correct
func CheckBucketClassSpec(clientset *kubernetes.Clientset, bucketClassSpec v1alpha1.BucketClaimSpec) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckSecret Check if secret exists
func CheckSecret(clientset *kubernetes.Clientset, secretName string, namespace string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketClaimEvents Check BucketClaim events
func CheckBucketClaimEvents(clientset *kubernetes.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}
