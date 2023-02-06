package steps

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// CheckCOSIControllerInstallation Ensure that COSI controller objectstorage-controller is installed in namespace "driver-ns"
func CheckCOSIControllerInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, controllerName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.HaveField("Type", gomega.Equal(v1.DeploymentAvailable))))
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in namespace "driver-ns"
func CheckCOSIDriverInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, driver string, namespace string) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, driver, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.HaveField("Type", gomega.Equal(v1.DeploymentAvailable))))
}
