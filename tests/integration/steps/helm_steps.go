package steps

import (
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// CheckCOSIControllerInstallation Ensure that COSI controller 'cosi-controller' is installed in namespace "driver-ns"
func CheckCOSIControllerInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, controllerName, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodRunning))
}

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet) {
	_, err := objectscale.FederatedObjectStores().List(map[string]string{})
	gomega.Expect(err).To(gomega.BeNil())
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in namespace "driver-ns"
func CheckCOSIDriverInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, driver string, namespace string) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, driver, metav1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodRunning))
}
