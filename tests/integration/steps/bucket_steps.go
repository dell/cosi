package steps

import (
	. "github.com/onsi/ginkgo/v2"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// Function creating a BucketClaim resource from specification
func CreateBucketClaimResource(bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Function for checking BucketClaim status
func CheckBucketClaimStatus(bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")

}

// Function for checking Bucket status and bucketID
func CheckBucketStatusAndBucketID(bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}

// Function for creating BucketClass resource
func CreateBucketClassResource(bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	// TODO: Implementation goes here
	Fail("UNIMPLEMENTED")
}
