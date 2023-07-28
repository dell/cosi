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
	"errors"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/dell/goobjectscale/pkg/client/model"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	objscl "github.com/dell/cosi/pkg/provisioner/objectscale"
	objectscaleRest "github.com/dell/goobjectscale/pkg/client/rest"
	gomega "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"

	"github.com/dell/cosi/pkg/provisioner/policy"
)

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster.
func CheckObjectScaleInstallation(ctx context.Context, objectscale *objectscaleRest.ClientSet, namespace string) {
	param := make(map[string]string)
	param["namespace"] = namespace

	_, err := objectscale.Buckets().List(ctx, param)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckObjectStoreExists Ensure that ObjectStore "${objectstoreId}" is created.
func CheckObjectStoreExists(ctx context.Context, objectscale *objectscaleRest.ClientSet, objectstore string) {
	objectstores, err := objectscale.FederatedObjectStores().List(ctx, make(map[string]string))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectstores.Items).To(gomega.ContainElement(gomega.HaveField("ObjectStoreID", objectstore)))
}

// CheckBucketResourceInObjectStore Function checking if Bucket resource is in objectstore.
func CheckBucketResourceInObjectStore(ctx context.Context, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = namespace
	id := strings.SplitN(bucket.Status.BucketID, "-", 2)[1] //nolint:gomnd

	objectScaleBucket, err := objectscale.Buckets().Get(ctx, id, param)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(objectScaleBucket).NotTo(gomega.BeNil())
}

// CheckBucketDeletionInObjectStore Function for checking Bucket deletion in ObjectStore.
func CheckBucketDeletionInObjectStore(ctx context.Context, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	param := make(map[string]string)
	param["namespace"] = namespace
	id := strings.SplitN(bucket.Status.BucketID, "-", 2)[1] //nolint:gomnd

	err := retry(ctx, attempts, sleep, func() error {
		var err error

		expectedError := model.Error{
			Code: model.CodeParameterNotFound,
		}

		_, err = objectscale.Buckets().Get(ctx, id, param)
		if !errors.Is(err, expectedError) {
			return err
		}

		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreatePolicy Function for creating policy in ObjectScale.
func CreatePolicy(ctx context.Context, objectscale *objectscaleRest.ClientSet, policy policy.Document, myBucket *v1alpha1.Bucket) {
	policyString, err := policy.ToJSON()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	err = objectscale.Buckets().UpdatePolicy(ctx, myBucket.Name, policyString, nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckPolicy checks  if policy exists in ObjectScale.
func CheckPolicy(ctx context.Context, objectscale *objectscaleRest.ClientSet, expectedPolicyDocument policy.Document, myBucket *v1alpha1.Bucket, namespace string) {
	var actualPolicyDocument policy.Document

	ErrComparisonFailed := errors.New("comparison failed")

	// This also needs to be retried, as we are not sure, if the policy was already updated.
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		param := make(map[string]string)
		param["namespace"] = namespace

		actualPolicy, err := objectscale.Buckets().GetPolicy(ctx, myBucket.Name, param)
		if err != nil {
			return err
		}

		actualPolicyDocument, err = policy.NewFromJSON(actualPolicy)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(actualPolicyDocument, expectedPolicyDocument) {
			return ErrComparisonFailed
		}

		return nil
	})

	// If the error is ErrComparisonFailed, I want full gomega match on objects, so I get pretty output.
	gomega.Expect(err).To(gomega.Or(gomega.BeEquivalentTo(ErrComparisonFailed), gomega.Not(gomega.HaveOccurred())))
	gomega.Expect(expectedPolicyDocument).To(gomega.BeEquivalentTo(actualPolicyDocument))
}

// DeletePolicy is a function deleting a policy from the ObjectStore.
func DeletePolicy(ctx context.Context, objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket, namespace string) {
	param := make(map[string]string)
	param["namespace"] = namespace
	existing, err := objectscale.Buckets().GetPolicy(ctx, bucket.Name, param)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(existing).NotTo(gomega.BeNil())

	err = objectscale.Buckets().DeletePolicy(ctx, bucket.Name, param)

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateUser creates user in ObjectScale.
func CreateUser(ctx context.Context, iamClient *iam.IAM, user string, arn string) {
	// TODO: verify it's working correctly once all the steps are integrated
	userOut, err := iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName:            &user,
		PermissionsBoundary: &arn,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
}

// CheckUser checks if user exists in ObjectScale.
func CheckUser(ctx context.Context, iamClient *iam.IAM, user string, namespace string) {
	username := objscl.BuildUsername(namespace, user)
	userOut, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &username})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
	gomega.Expect(username).To(gomega.Equal(*(userOut.User.UserName)))
}

// CheckUserDeleted checks if user does not exist in ObjectScale.
func CheckUserDeleted(ctx context.Context, iamClient *iam.IAM, user string, namespace string) {
	username := objscl.BuildUsername(namespace, user)

	_, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &username})
	gomega.Expect(err).To(gomega.HaveOccurred())

	var myAwsErr awserr.Error
	matched := errors.As(err, &myAwsErr)

	gomega.Expect(matched).To(gomega.BeTrue())
	gomega.Expect(iam.ErrCodeNoSuchEntityException).To(gomega.Equal(myAwsErr.Code()))
}

// DeleteUser Function for deleting user from ObjectScale.
func DeleteUser(ctx context.Context, iamClient *iam.IAM, user string) {
	existing, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &user})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(existing.User).NotTo(gomega.BeNil())
	// TODO: verify it's working correctly once all the steps are integrated
	_, err = iamClient.DeleteUser(&iam.DeleteUserInput{UserName: existing.User.UserName})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketNotInObjectStore Function for checking if bucket is not in objectstore.
func CheckBucketNotInObjectStore(ctx context.Context, objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(ctx, bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).To(gomega.HaveOccurred())
	gomega.Expect(bucket).To(gomega.BeNil())
}

// CheckBucketInObjectStore Function for checking if the bucket object is in the objectstore.
func CheckBucketInObjectStore(ctx context.Context, objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(ctx, bucketClaim.Status.BucketName, map[string]string{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())
}

// DeleteBucket Function for deleting existing from ObjectScale (useful if BucketClaim deletionPolicy is set to "retain").
func DeleteBucket(ctx context.Context, objectscale *objectscaleRest.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	err := objectscale.Buckets().Delete(ctx, bucket.Name, namespace, false)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

func DeleteAccessKey(ctx context.Context, iamClient *iam.IAM, clientset *kubernetes.Clientset, secret *v1.Secret) {
	accessKeyID := GetAccessKeyID(ctx, clientset, secret)
	_, err := iamClient.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
		AccessKeyId: &accessKeyID,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}
