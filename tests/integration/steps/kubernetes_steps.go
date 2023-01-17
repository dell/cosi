package steps

import (
	. "github.com/onsi/ginkgo/v2"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// Ensure that Kubernetes cluster is available
func CheckClusterAvailability(clientset *kubernetes.Clientset) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Ensure that Kubernetes namespace "driver-ns" is created
func CreateNamespace(clientset *kubernetes.Clientset, namespace string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Ensure that specification of custom resource "my-bucket-class" is correct
func CheckBucketClassSpec(clientset *kubernetes.Clientset, bucketClassSpec v1alpha1.BucketClaimSpec) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

func CheckSecret(clientset *kubernetes.Clientset, secretName string, namespace string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}
