package steps

import (
	"net/http"
	"path"

	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest/client"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// CheckCOSIControllerInstallation Ensure that COSI controller 'cosi-controller' is installed in namespace "driver-ns"
func CheckCOSIControllerInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(len(pods.Items)).To(gomega.BeNumerically(">", 0))
	pod := pods.Items[0]
	for p := range pods.Items {
		if pods.Items[p].Name == controllerName {
			pod = pods.Items[p]
			break
		}
	}
	gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodRunning))
}

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet) {
	req := client.Request{
		Method:      http.MethodGet,
		Path:        path.Join("fedsvc", "objectScaleName"),
		ContentType: client.ContentTypeXML,
	}
	var objectscaleName string
	err := objectscale.Client().MakeRemoteCall(req, &objectscaleName)
	gomega.Expect(err).To(gomega.BeNil())
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in namespace "driver-ns"
func CheckCOSIDriverInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, driver string, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(len(pods.Items)).To(gomega.BeNumerically(">", 0))
	pod := pods.Items[0]
	for p := range pods.Items {
		if pods.Items[p].Name == driver {
			pod = pods.Items[p]
			break
		}
	}
	gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodRunning))
}
