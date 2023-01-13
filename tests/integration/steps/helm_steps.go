package steps

import (
	. "github.com/onsi/ginkgo/v2"
	"k8s.io/client-go/kubernetes"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// Ensure that COSI controller is installed in namespace "driver-ns"
func CheckCOSIControllerInstallation(clientset *kubernetes.Clientset, namespace string) {
	// TODO: Implementation goes here
	// check if COSI controller is installed in namespace "driver-ns"
	// if not, fail the test
	Fail("UNIMPLEMENTED")
}

// Ensure that ObjectScale platform is installed on the cluster
func CheckObjectScaleInstallation(clientset *kubernetes.Clientset) {
	// TODO: Implementation goes here
	// check if ObjectScale is installed
	// if not, fail the test
	Fail("UNIMPLEMENTED")
}

// Ensure that COSI driver is installed in namespace "driver-ns"
func CheckCOSIDriverInstallation(clientset *kubernetes.Clientset, driver string, namespace string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}
