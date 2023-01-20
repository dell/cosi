package steps

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// CreateBucketClaimResource Function creating a BucketClaim resource from specification
func CreateBucketClaimResource(bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketClaimStatus Function for checking BucketClaim status
func CheckBucketClaimStatus(bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")

}

// CheckBucketStatus Function for checking Bucket status
func CheckBucketStatus(bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketID Function for checking bucketID
func CheckBucketID(bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreateBucketClassResource Function for creating BucketClass resource
func CreateBucketClassResource(bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// DeleteBucketClaimResource Function for deleting BucketClaim resource
func DeleteBucketClaimResource(bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreateBucketAccessClassResource Function for creating BucketAccessClass resource
func CreateBucketAccessClassResource(bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CreateBucketAccessResource Function for creating BucketAccess resource
func CreateBucketAccessResource(bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketAccessStatus Function for checking BucketAccess status
func CheckBucketAccessStatus(bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// CheckBucketAccessAccountID Function for checking BucketAccess accountID
func CheckBucketAccessAccountID(bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, accountID string) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}

// DeleteBucketAccessResource Function for deleting BucketAccess resource
func DeleteBucketAccessResource(bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	// TODO: Implementation goes here
	ginkgo.Fail("UNIMPLEMENTED")
}
