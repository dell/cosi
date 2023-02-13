package steps

import (
	"fmt"
	"log"
	"os"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed
)

// CheckCOSIControllerInstallation Ensure that COSI controller objectstorage-controller is installed in particular namespace
func CheckCOSIControllerInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	checkAppIsInstalled(ctx, clientset, controllerName, namespace)
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in particular namespace
func CheckCOSIDriverInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, driver string, namespace string) {
	checkAppIsInstalled(ctx, clientset, driver, namespace)
}

// checkAppIsInstalled Ensures that an app is installed in particular namespace
func checkAppIsInstalled(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, releaseName string, namespace string) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, releaseName, metav1.GetOptions{})
	if err != nil {
		InstallChartInNamespace(releaseName, namespace)
	} else {
		gomega.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.HaveField("Type", gomega.Equal(v1.DeploymentAvailable))))
	}
}

// InstallChart
func InstallChartInNamespace(releaseName, namespace string) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), log.Printf)
	gomega.Expect(err).To(gomega.BeNil())

	helmClient := action.NewInstall(actionConfig)
	helmClient.ReleaseName = releaseName
	helmClient.Namespace = namespace

	chartPath, err := helmClient.LocateChart("https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.0.6/ingress-nginx-4.0.6.tgz", settings)

	gomega.Expect(err).To(gomega.BeNil())
	chart, err := loader.Load(chartPath)
	gomega.Expect(err).To(gomega.BeNil())

	// panic(fmt.Sprintf("%v, %v", helmClient.ReleaseName, helmClient.Namespace))

	helmClient.DryRun = true
	release, err := helmClient.Run(chart, nil)
	gomega.Expect(err).To(gomega.BeNil())

	fmt.Println("Successfully installed release: ", release.Name)
}

// UninstallChartReleaseinNamespace Deletes particular relese from k8s chart
func UninstallChartReleaseinNamespace(releaseName string, namespace string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}
