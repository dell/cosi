package steps

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"k8s.io/client-go/kubernetes"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// CheckCOSIControllerInstallation Ensure that COSI controller 'cosi-controller' is installed in namespace "driver-ns"
func CheckCOSIControllerInstallation(clientset *kubernetes.Clientset, controllerName string, namespace string) {
	// TODO: Implementation goes here
	// check if COSI controller is installed in namespace "driver-ns"
	// if not, fail the test
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset) {
	// TODO: Implementation goes here
	// check if ObjectScale is installed
	// if not, fail the test
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in namespace "driver-ns"
func CheckCOSIDriverInstallation(clientset *kubernetes.Clientset, driver string, namespace string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}
