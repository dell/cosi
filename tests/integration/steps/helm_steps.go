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
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	// TODO: use https://pkg.go.dev/helm.sh/helm/v3 for helm operations if needed.
)

// CheckCOSIControllerInstallation Ensure that COSI controller objectstorage-controller is installed in particular namespace.
func CheckCOSIControllerInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	// TODO: check if controller can be installed via chart
	repo := ""
	chartName := ""
	version := ""
	checkAppIsInstalled(ctx, clientset, controllerName, namespace, repo, chartName, version)
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in particular namespace.
func CheckCOSIDriverInstallation(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, driver string, namespace string) {
	repo := "cosi-driver"
	chartName := "cosi-driver"
	version := "0.1.0"
	checkAppIsInstalled(ctx, clientset, driver, namespace, repo, chartName, version)
}

// checkAppIsInstalled Ensure that an app is installed in particular namespace.
func checkAppIsInstalled(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, releaseName, namespace, repo, chartName, version string) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, releaseName, metav1.GetOptions{})
	if err != nil {
		InstallChartInNamespace(releaseName, namespace, repo, chartName, version)
	} else {
		gomega.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.HaveField("Type", gomega.Equal(v1.DeploymentAvailable))))
	}
}

// InstallChart Install particular release from k8s chart.
func InstallChartInNamespace(releaseName, namespace, repo, chartName, version string) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Debugf)
	gomega.Expect(err).To(gomega.BeNil())

	helmClient := action.NewInstall(actionConfig)
	helmClient.ReleaseName = releaseName
	helmClient.Namespace = namespace

	chartPath, err := helmClient.LocateChart(fmt.Sprintf("https://github.com/%s/%s-%s", repo, chartName, version), settings)

	gomega.Expect(err).To(gomega.BeNil())
	chart, err := loader.Load(chartPath)
	gomega.Expect(err).To(gomega.BeNil())

	release, err := helmClient.Run(chart, nil)
	gomega.Expect(err).To(gomega.BeNil())

	log.Println("Successfully installed release: ", release.Name)
}

// UninstallChartReleaseinNamespace Delete particular release from k8s chart.
func UninstallChartReleaseinNamespace(releaseName, namespace string) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Debugf)
	gomega.Expect(err).To(gomega.BeNil())

	helmClient := action.NewUninstall(actionConfig)
	release, err := helmClient.Run(releaseName)
	gomega.Expect(err).To(gomega.BeNil())

	log.Println("Successfully uninstalled release: ", release.Release.Name)
}
