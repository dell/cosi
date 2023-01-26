package steps

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster
func CheckObjectScaleInstallation(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet) {
	_, err := objectscale.FederatedObjectStores().List(map[string]string{})
	gomega.Expect(err).To(gomega.BeNil())
}

// CheckObjectStoreCreation Ensure that ObjectStore "objectstore-dev" is created
func CheckObjectStoreExists(ctx ginkgo.SpecContext, objectscale *objectscaleRest.ClientSet, objectstore string) {
	objectstores, err := objectscale.FederatedObjectStores().List(make(map[string]string))
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(objectstores.Items).To(gomega.ContainElement(gomega.HaveField("ObjectStoreName", objectstore)))
}

// CheckBucketResourceInObjectStore Function checking if Bucket resource is in objectstore
func CheckBucketResourceInObjectStore(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketDeletionInObjectStore Function for checking Bucket deletion in ObjectStore
func CheckBucketDeletionInObjectStore(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketAccessFromSecret Check if Bucket can be accessed with data from specified secret
func CheckBucketAccessFromSecret(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket, secretName string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreatePolicy Function for creating policy in ObjectScale
func CreatePolicy(objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	err := objectscale.Buckets().UpdatePolicy(myBucket.Name, policy, nil)
	gomega.Expect(err).To(gomega.BeNil())
	ginkgo.Fail("UNIMPLEMENTED")
}

// Function for checking if policy exists in ObjectScale
// TODO: responisbility of @shanduur-dell
func CheckPolicy(objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	actualPolicy, err := objectscale.Buckets().GetPolicy(myBucket.Name, nil)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(actualPolicy).To(gomega.BeIdenticalTo(policy))
}

// Function for creating user in ObjectScale
func CreateUser(ctx context.Context, iamClient *iam.IAM, user, arn string) {
	// TODO: Implementation goes here
	userOut, err := iamClient.CreateUserWithContext(ctx, &iam.CreateUserInput{
		UserName:            &user,
		PermissionsBoundary: &arn,
	})
	if gomega.Expect(err).To(gomega.BeNil()) {
		gomega.Expect(userOut.User).NotTo(gomega.BeNil())
	}
	ginkgo.Fail("UNIMPLEMENTED")
}

// Function for checking if user exists in ObjectScale
// ASSIGNEE: @shanduur-dell
func CheckUser(ctx context.Context, iamClient *iam.IAM, user string) {
	userOut, err := iamClient.GetUserWithContext(ctx, &iam.GetUserInput{UserName: &user})
	if gomega.Expect(err).To(gomega.BeNil()) {
		gomega.Expect(userOut.User).NotTo(gomega.BeNil())
		gomega.Expect(*(userOut.User.UserName)).To(gomega.Equal(user))
		gomega.Expect(userOut.User.Arn).To(gomega.Or(gomega.BeNil(), gomega.BeEmpty()))
	}
}

// DeletePolicy Function deleteing policy from ObjectStore
func DeletePolicy(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// DeleteUser Function for deleting user from ObjectScale
func DeleteUser(objectscale *objectscaleRest.ClientSet, user string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketNotInObjectStore Function for checking if bucket is not in objectstore
// ASSIGNEE: @shanduur-dell
func CheckBucketNotInObjectStore(objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	bucket, err := objectscale.Buckets().Get(bucketClaim.Name, map[string]string{})
	gomega.Expect(err).NotTo(gomega.BeNil())
	gomega.Expect(bucket).To(gomega.BeNil())
}
