// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package steps

import (
	"context"
	"fmt"

	"sigs.k8s.io/container-object-storage-interface/client/apis/objectstorage/v1alpha1"

	gomega "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface/client/clientset/versioned"
)

// CreateBucketClaimResource Function creating a BucketClaim resource from specification.
func CreateBucketResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	_, err := bucketClient.ObjectstorageV1alpha1().Buckets().Create(ctx, bucket, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucketClaimResource Function creating a BucketClaim resource from specification.
func CreateBucketClaimResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Create(ctx, bucketClaim, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// DeleteBucketClaimResource Function for deleting BucketClaim resource.
func DeleteBucketClaimResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Delete(ctx, bucketClaim.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketClaimStatus Function for checking BucketClaim status.
func CheckBucketClaimStatus(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim, status bool) {
	var myBucketClaim *v1alpha1.BucketClaim

	claimAttempt := 0
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		claimAttempt++

		myBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
		if err != nil {
			fmt.Printf("[CheckBucketClaimStatus] attempt %d/%d: error fetching BucketClaim %s/%s: %v\n", claimAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name, err)
			return err
		}

		if myBucketClaim.Status.BucketReady != status {
			fmt.Printf("[CheckBucketClaimStatus] attempt %d/%d: BucketReady is %v, expected %v for %s/%s\n", claimAttempt, attempts, myBucketClaim.Status.BucketReady, status, bucketClaim.Namespace, bucketClaim.Name)
			return fmt.Errorf("BucketReady is %v, expected %v", myBucketClaim.Status.BucketReady, status)
		}

		fmt.Printf("[CheckBucketClaimStatus] attempt %d/%d: BucketClaim %s/%s has expected status %v\n", claimAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name, status)
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim).NotTo(gomega.BeNil())
	gomega.Expect(myBucketClaim.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketStatus Function for checking Bucket status.
func CheckBucketStatus(bucket *v1alpha1.Bucket, status bool) {
	gomega.Expect(bucket.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketID Function for checking bucketID.
func CheckBucketID(bucket *v1alpha1.Bucket) {
	gomega.Expect(bucket.Status.BucketID).NotTo(gomega.Or(gomega.BeEmpty(), gomega.BeNil()))
}

// CreateBucketClassResource Function for creating BucketClass resource.
func CreateBucketClassResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) *v1alpha1.BucketClass {
	bucketClass, err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Create(ctx, bucketClass, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	return bucketClass
}

// DeleteBucketClassResource Function for deleting BucketClass resource.
func DeleteBucketClassResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Delete(ctx, bucketClass.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucketAccessClassResource Function for creating BucketAccessClass resource.
func CreateBucketAccessClassResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Create(ctx, bucketAccessClass, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// DeleteBucketAccessClassResource Function for deleting BucketAccessClass resource.
func DeleteBucketAccessClassResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Delete(ctx, bucketAccessClass.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucketAccessResource Function for creating BucketAccess resource.
func CreateBucketAccessResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Create(ctx, bucketAccess, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// DeleteBucketAccessResource Function for deleting BucketAccess resource.
func DeleteBucketAccessResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	// Not sure why but ginkgo context is messing with kubernetes client here creating error:
	// client rate limiter Wait returned error: context canceled
	// So new context it is.
	err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Delete(ctx, bucketAccess.Name, v1.DeleteOptions{})
	// it's ok if it's no error or erorr has reason NotFound, we want to delete it anyway
	gomega.Expect(err).To(gomega.Or(gomega.BeNil(), gomega.HaveField("Reason", "NotFound")))
}

// CheckBucketAccessStatus Function for checking BucketAccess status.
func CheckBucketAccessStatus(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, status bool) *v1alpha1.BucketAccess {
	var myBucketAccess *v1alpha1.BucketAccess

	accessAttempt := 0
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		accessAttempt++

		myBucketAccess, err = bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
		if err != nil {
			fmt.Printf("[CheckBucketAccessStatus] attempt %d/%d: error fetching BucketAccess %s/%s: %v\n", accessAttempt, attempts, bucketAccess.Namespace, bucketAccess.Name, err)
			return err
		}

		if !myBucketAccess.Status.AccessGranted {
			fmt.Printf("[CheckBucketAccessStatus] attempt %d/%d: AccessGranted is false for %s/%s\n", accessAttempt, attempts, bucketAccess.Namespace, bucketAccess.Name)
			return fmt.Errorf("AccessGranted is false")
		}

		fmt.Printf("[CheckBucketAccessStatus] attempt %d/%d: BucketAccess %s/%s has AccessGranted=%v\n", accessAttempt, attempts, bucketAccess.Namespace, bucketAccess.Name, myBucketAccess.Status.AccessGranted)
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(myBucketAccess.Status.AccessGranted).To(gomega.Equal(status))

	return myBucketAccess
}

// CheckBucketAccessAccountID Function for checking BucketAccess accountID.
func CheckBucketAccessAccountID(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, accountID string) {
	myBucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(accountID).To(gomega.Equal(myBucketAccess.Status.AccountID))
}

// GetBucketResource function for getting Bucket resource.
func GetBucketResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) *v1alpha1.Bucket {
	var myBucketClaim *v1alpha1.BucketClaim

	claimAttempt := 0
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		claimAttempt++

		myBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
		if err != nil {
			fmt.Printf("[GetBucketResource] attempt %d/%d: error fetching BucketClaim %s/%s: %v\n", claimAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name, err)
			return err
		}

		if myBucketClaim.Status.BucketName == "" && myBucketClaim.Spec.ExistingBucketName == "" {
			fmt.Printf("[GetBucketResource] attempt %d/%d: BucketClaim %s/%s still missing bucket reference\n", claimAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name)
			return fmt.Errorf("BucketName and ExistingBucketName are empty")
		}

		fmt.Printf("[GetBucketResource] attempt %d/%d: BucketClaim %s/%s has bucket reference\n", claimAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name)
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	var bucket *v1alpha1.Bucket

	bucketAttempt := 0
	err = retry(ctx, attempts, sleep, func() error {
		var err error

		bucketAttempt++

		name := ""
		if myBucketClaim.Spec.ExistingBucketName != "" {
			name = myBucketClaim.Spec.ExistingBucketName
		} else {
			name = myBucketClaim.Status.BucketName
		}

		bucket, err = bucketClient.ObjectstorageV1alpha1().Buckets().Get(ctx, name, v1.GetOptions{})
		if err != nil {
			fmt.Printf("[GetBucketResource] attempt %d/%d: error fetching Bucket %s: %v\n", bucketAttempt, attempts, name, err)
			return err
		}

		if bucket.Spec.ExistingBucketID != "" {
			if bucket.Status.BucketReady && bucket.Status.BucketID != "" {
				fmt.Printf("[GetBucketResource] attempt %d/%d: bucket %s is ready\n", bucketAttempt, attempts, bucket.Name)
				return nil
			}

			fmt.Printf("[GetBucketResource] attempt %d/%d: bucket %s is not ready yet\n", bucketAttempt, attempts, bucket.Name)
			return fmt.Errorf("bucket %s is not ready yet", bucket.Name)
		}

		if !bucket.Status.BucketReady {
			fmt.Printf("[GetBucketResource] attempt %d/%d: bucket %s is not ready\n", bucketAttempt, attempts, bucket.Name)
			return fmt.Errorf("bucket %s is not ready", bucket.Name)
		}

		fmt.Printf("[GetBucketResource] attempt %d/%d: bucket %s is ready\n", bucketAttempt, attempts, bucket.Name)
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())

	return bucket
}

// CheckBucketStatusEmpty function for checking if Bucket status is empty.
func CheckBucketStatusEmpty(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	var myBucketClaim *v1alpha1.BucketClaim

	emptyAttempt := 0
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		emptyAttempt++

		myBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
		if err != nil {
			fmt.Printf("[CheckBucketStatusEmpty] attempt %d/%d: error fetching BucketClaim %s/%s: %v\n", emptyAttempt, attempts, bucketClaim.Namespace, bucketClaim.Name, err)
		}

		return err
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim.Status.BucketName).To(gomega.BeEmpty())
}
