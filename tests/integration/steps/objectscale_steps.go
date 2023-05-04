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

	"github.com/aws/aws-sdk-go/service/iam"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
)

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster.
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet) {
	_, err := objectscale.FederatedObjectStores().List(map[string]string{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckObjectStoreCreation Ensure that ObjectStore "${objectstoreName}" is created.
func CheckObjectStoreExists(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, objectstore string) {
	objectstores, err := objectscale.FederatedObjectStores().List(make(map[string]string))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectstores.Items).To(gomega.ContainElement(gomega.HaveField("ObjectStoreName", objectstore)))
}

// CheckBucketResourceInObjectStore Function checking if Bucket resource is in objectstore.
func CheckBucketResourceInObjectStore(objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = namespace

	id := strings.SplitN(bucket.Status.BucketID, "-", 2)[1]

	objectScaleBucket, err := objectscale.Buckets().Get(id, param)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectScaleBucket).NotTo(gomega.BeNil())
}

// CheckBucketDeletionInObjectStore Function for checking Bucket deletion in ObjectStore.
func CheckBucketDeletionInObjectStore(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = "TODO:Separate-ObjectStoreID-from-bucket"
	objectScaleBucket, err := objectscale.Buckets().Get(bucket.Status.BucketID, param)
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(objectScaleBucket).To(gomega.BeNil())
}

// CheckBucketAccessFromSecret Check if Bucket can be accessed with data from specified secret.
func CheckBucketAccessFromSecret(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket, secretName string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreatePolicy Function for creating policy in ObjectScale.
func CreatePolicy(objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	err := objectscale.Buckets().UpdatePolicy(myBucket.Name, policy, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// Function for checking if policy exists in ObjectScale.
func CheckPolicy(objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	actualPolicy, err := objectscale.Buckets().GetPolicy(myBucket.Name, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(actualPolicy).To(gomega.BeIdenticalTo(policy))
}

// DeletePolicy is a function deleting a policy from the ObjectStore.
func DeletePolicy(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	existing, err := objectscale.Buckets().GetPolicy(bucket.Name, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(existing).NotTo(gomega.BeNil())

	err = objectscale.Buckets().DeletePolicy(bucket.Name, nil)

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// Function for creating user in ObjectScale.
func CreateUser(ctx ginkgo.SpecContext, iamClient *iam.IAM, user, arn string) {
	// TODO: verify it's working correctly once all the steps are integrated
	userOut, err := iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName:            &user,
		PermissionsBoundary: &arn,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
}

// Function for checking if user exists in ObjectScale.
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
func CheckBucketNotInObjectStore(objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(bucket).To(gomega.BeNil())
}

// CheckBucketInObjectStore Function for checking if the bucket object is in the objectstore.
func CheckBucketInObjectStore(objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())
}

// DeleteBucket Function for deleting existing from ObjectScale (useful if BucketClaim deletionPolicy is set to "retain").
func DeleteBucket(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	err := objectscale.Buckets().Delete(bucket.Name, bucket.Namespace, true)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}
