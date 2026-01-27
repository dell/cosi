// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package steps

import (
	"context"

	"k8s.io/client-go/kubernetes"

	gomega "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CheckCOSIControllerInstallation Ensure that COSI controller objectstorage-controller is installed in particular namespace.
func CheckCOSIControllerInstallation(ctx context.Context, clientset *kubernetes.Clientset, controllerName string, namespace string) {
	// TODO: check if controller can be installed via chart
	checkAppIsInstalled(ctx, clientset, controllerName, namespace)
}

// CheckCOSIDriverInstallation Ensure that COSI driver is installed in particular namespace.
func CheckCOSIDriverInstallation(ctx context.Context, clientset *kubernetes.Clientset, driver string, namespace string) {
	checkAppIsInstalled(ctx, clientset, driver, namespace)
}

// checkAppIsInstalled Ensure that an app is installed in particular namespace.
func checkAppIsInstalled(ctx context.Context, clientset *kubernetes.Clientset, releaseName, namespace string) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, releaseName, metav1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(deployment.Status.Conditions).To(gomega.ContainElement(gomega.HaveField("Type", gomega.Equal(v1.DeploymentAvailable))))
}
