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
	"strings"

	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/service/iam"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster.
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet) {
	_, err := objectscale.FederatedObjectStores().List(ctx, map[string]string{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckObjectStoreExists Ensure that ObjectStore "${objectstoreId}" is created.
func CheckObjectStoreExists(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, objectstore string) {
	objectstores, err := objectscale.FederatedObjectStores().List(ctx, make(map[string]string))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectstores.Items).To(gomega.ContainElement(gomega.HaveField("ObjectStoreID", objectstore)))
}

// CheckBucketResourceInObjectStore Function checking if Bucket resource is in objectstore.
func CheckBucketResourceInObjectStore(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = namespace
	id := strings.SplitN(bucket.Status.BucketID, "-", 2)[1] // nolint:gomnd

	objectScaleBucket, err := objectscale.Buckets().Get(ctx, id, param)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectScaleBucket).NotTo(gomega.BeNil())

	ginkgo.GinkgoWriter.Printf("Bucket in Objectstore: %+v\n", objectScaleBucket)
}

// CheckBucketDeletionInObjectStore Function for checking Bucket deletion in ObjectStore.
func CheckBucketDeletionInObjectStore(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = namespace
	id := strings.SplitN(bucket.Status.BucketID, "-", 2)[1] // nolint:gomnd

	err := retry(ctx, attempts, sleep, func() error {
		var err error
		objectScaleBucket, err := objectscale.Buckets().Get(ctx, id, param)
		ginkgo.GinkgoWriter.Printf("Bucket in Objectstore: %+v\n", objectScaleBucket)
		return err
	})

	gomega.Expect(err).To(gomega.HaveOccurred())
}

// CheckBucketAccessFromSecret Check if Bucket can be accessed with data from specified secret.
func CheckBucketAccessFromSecret(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket, secretName string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreatePolicy Function for creating policy in ObjectScale.
func CreatePolicy(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	err := objectscale.Buckets().UpdatePolicy(ctx, myBucket.Name, policy, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckPolicy checks  if policy exists in ObjectScale.
func CheckPolicy(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	actualPolicy, err := objectscale.Buckets().GetPolicy(ctx, myBucket.Name, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(actualPolicy).To(gomega.BeIdenticalTo(policy))
}

// DeletePolicy is a function deleting a policy from the ObjectStore.
func DeletePolicy(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	existing, err := objectscale.Buckets().GetPolicy(ctx, bucket.Name, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(existing).NotTo(gomega.BeNil())

	err = objectscale.Buckets().DeletePolicy(ctx, bucket.Name, nil)

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateUser creates user in ObjectScale.
func CreateUser(ctx ginkgo.SpecContext, iamClient *iam.IAM, user, arn string) {
	// TODO: verify it's working correctly once all the steps are integrated
	userOut, err := iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName:            &user,
		PermissionsBoundary: &arn,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
}

// CheckUser checks if user exists in ObjectScale.
func CheckUser(ctx ginkgo.SpecContext, iamClient *iam.IAM, user string) {
	// TODO: verify it's working correctly once all the steps are integrated
	userOut, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &user})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
	gomega.Expect(*(userOut.User.UserName)).To(gomega.Equal(user))
	gomega.Expect(userOut.User.Arn).To(gomega.Or(gomega.BeNil(), gomega.BeEmpty()))
}

// DeleteUser Function for deleting user from ObjectScale.
func DeleteUser(ctx ginkgo.SpecContext, iamClient *iam.IAM, user string) {
	existing, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &user})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(existing.User).NotTo(gomega.BeNil())
	// TODO: verify it's working correctly once all the steps are integrated
	_, err = iamClient.DeleteUser(&iam.DeleteUserInput{UserName: existing.User.UserName})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketNotInObjectStore Function for checking if bucket is not in objectstore.
func CheckBucketNotInObjectStore(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(ctx, bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(bucket).To(gomega.BeNil())
}

// CheckBucketInObjectStore Function for checking if the bucket object is in the objectstore.
func CheckBucketInObjectStore(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(ctx, bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())
}

// DeleteBucket Function for deleting existing from ObjectScale (useful if BucketClaim deletionPolicy is set to "retain").
func DeleteBucket(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	err := objectscale.Buckets().Delete(ctx, bucket.Name, namespace, false)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}
