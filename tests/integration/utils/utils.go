package utils

import (
	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/dell/cosi-driver/tests/integration/steps"
	"k8s.io/client-go/kubernetes"
)

func DeleteReleasesAndNamespaces(ctx ginkgo.SpecContext, clientset *kubernetes.Clientset, releases map[string]string, namespaces []string) {
	// uninstall releases
	for namespace, release := range releases {
		steps.UninstallChartReleaseinNamespace(release, namespace)
	}
	// delete namespaces
	for _, namespace := range namespaces {
		steps.DeleteNamespace(ctx, clientset, namespace)
	}
}
