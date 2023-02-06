package steps

import (
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"
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
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreateUser Function for creating user in ObjectScale
func CreateUser(objectscale *objectscaleRest.ClientSet, user string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
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
func CheckBucketNotInObjectStore(objectscale *objectscaleRest.ClientSet, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}
