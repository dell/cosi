package steps

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	errors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// CreateBucketClaimResource Function creating a BucketClaim resource from specification
func CreateBucketClaimResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Create(ctx, bucketClaim, v1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		gomega.Expect(err).To(gomega.BeNil())
	}
}

// CheckBucketClaimStatus Function for checking BucketClaim status
func CheckBucketClaimStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim, status bool) {
	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(myBucketClaim).NotTo(gomega.BeNil())
	gomega.Expect(myBucketClaim.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketStatus Function for checking Bucket status
func CheckBucketStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket, status bool) {
	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucket.Spec.BucketClaim.Namespace).Get(ctx, bucket.Spec.BucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	myBucket, err := bucketClient.ObjectstorageV1alpha1().Buckets().Get(ctx, myBucketClaim.Status.BucketName, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(myBucket).NotTo(gomega.BeNil())
	gomega.Expect(myBucket.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketID Function for checking bucketID
func CheckBucketID(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucket.Spec.BucketClaim.Namespace).Get(ctx, bucket.Spec.BucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	myBucket, err := bucketClient.ObjectstorageV1alpha1().Buckets().Get(ctx, myBucketClaim.Status.BucketName, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(myBucket).NotTo(gomega.BeNil())
	gomega.Expect(myBucket.Status.BucketID).NotTo(gomega.Or(gomega.BeEmpty(), gomega.BeNil()))
}

// CreateBucketClassResource Function for creating BucketClass resource
func CreateBucketClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Get(ctx, bucketClass.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Create(ctx, bucketClass, v1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		gomega.Expect(err).To(gomega.BeNil())
	}
}

// DeleteBucketClassResource Function for deleting BucketClass resource
func DeleteBucketClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Delete(ctx, bucketClass.Name, v1.DeleteOptions{})
	gomega.Expect(err).To(gomega.BeNil())
}

// DeleteBucketClaimResource Function for deleting BucketClaim resource
func DeleteBucketClaimResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Delete(ctx, bucketClaim.Name, v1.DeleteOptions{})
	gomega.Expect(err).To(gomega.BeNil())
}

// CreateBucketAccessClassResource Function for creating BucketAccessClass resource
func CreateBucketAccessClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Get(ctx, bucketAccessClass.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Create(ctx, bucketAccessClass, v1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		gomega.Expect(err).To(gomega.BeNil())
	}
}

// CreateBucketAccessResource Function for creating BucketAccess resource
func CreateBucketAccessResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Create(ctx, bucketAccess, v1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		gomega.Expect(err).To(gomega.BeNil())
	}
}

// CheckBucketAccessStatus Function for checking BucketAccess status
func CheckBucketAccessStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, status bool) {
	myBucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(myBucketAccess.Status.AccessGranted).To(gomega.Equal(status))
}

// CheckBucketAccessAccountID Function for checking BucketAccess accountID
func CheckBucketAccessAccountID(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, accountID string) {
	myBucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(myBucketAccess.Status.AccountID).To(gomega.Equal(accountID))
}

// DeleteBucketAccessResource Function for deleting BucketAccess resource
func DeleteBucketAccessResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Delete(ctx, bucketAccess.Name, v1.DeleteOptions{})
	gomega.Expect(err).To(gomega.BeNil())
}
