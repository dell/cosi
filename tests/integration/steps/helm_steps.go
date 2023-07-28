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
