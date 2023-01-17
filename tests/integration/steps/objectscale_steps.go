package steps

import (
	objectscaleRest "github.com/emcecs/objectscale-management-go-sdk/pkg/client/rest"
	. "github.com/onsi/ginkgo/v2"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
)

// Ensure that ObjectStore "object-store-1" is created
func CheckObjectStoreCreation(objectscale *objectscaleRest.ClientSet, objectstore string) {
	// TODO: Implementation goes here
	// check if ObjectStore "object-store-1" is created
	// if not, fail the test
	Fail("UNIMPLEMENTED")
}

// Function checking if Bucket resource is in objectstore
func CheckBucketResourceInObjectStore(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Fucntrion for creating ObejctStore
func CreateObjectStore(objectscale *objectscaleRest.ClientSet, objectstore string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Function for checking Bucket deletion in ObjectStore
func CheckBucketDeletionInObjectStore(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

func CheckBucketAccessFromSecret(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket, secretName string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

func CreatePolicy(objectscale *objectscaleRest.ClientSet, policy string, myBucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

func CreateUser(objectscale *objectscaleRest.ClientSet, user string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Function deleteing policy from ObjectStore
func DeletePolicy(objectscale *objectscaleRest.ClientSet, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Function for deleting user from ObjectScale
func DeleteUser(objectscale *objectscaleRest.ClientSet, user string) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}
