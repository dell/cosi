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
	"fmt"

	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"

	gomega "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// CreateBucketClaimResource Function creating a BucketClaim resource from specification.
func CreateBucketClaimResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Create(ctx, bucketClaim, v1.CreateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// DeleteBucketClaimResource Function for deleting BucketClaim resource.
func DeleteBucketClaimResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// first delete finalizers
	// bucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	// gomega.Expect(err).ToNot(gomega.HaveOccurred())

	// bucketClaim.Finalizers = []string{}

	// _, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Update(ctx, bucketClaim, v1.UpdateOptions{})
	// gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Delete(ctx, bucketClaim.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketClaimStatus Function for checking BucketClaim status.
func CheckBucketClaimStatus(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim, status bool) {
	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim).NotTo(gomega.BeNil())
	gomega.Expect(myBucketClaim.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketStatus Function for checking Bucket status.
func CheckBucketStatus(ctx context.Context, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket, status bool) {
	gomega.Expect(bucket.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketID Function for checking bucketID.
func CheckBucketID(ctx context.Context, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
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
	// Not sure why but ginkgo context is meesing with kubernetes client here creating error:
	// client rate limiter Wait returned error: context canceled
	// So new context it is.
	bucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	// it's ok if it's no error or erorr has reason NotFound, we want to delete it anyway
	gomega.Expect(err).To(gomega.Or(gomega.BeNil(), gomega.HaveField("Reason", "NotFound")))
	// we can abort if bucketAccess is already deleted
	if err != nil {
		return
	}
	// remove finalizers deletes the BukcetAccess on cluster. I think controller does this but haven't checked.
	bucketAccess.Finalizers = []string{}
	_, err = bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Update(ctx, bucketAccess, v1.UpdateOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketAccessStatus Function for checking BucketAccess status.
func CheckBucketAccessStatus(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, status bool) *v1alpha1.BucketAccess {
	var myBucketAccess *v1alpha1.BucketAccess

	err := retry(ctx, attempts, sleep, func() error {
		var err error
		myBucketAccess, err = bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
		if err != nil {
			return err
		}

		if !myBucketAccess.Status.AccessGranted {
			return fmt.Errorf("AccessGranted is false")
		}
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
	gomega.Expect(myBucketAccess.Status.AccountID).To(gomega.Equal(accountID))
}

// GetBucketResource function for getting Bucket resource.
func GetBucketResource(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) *v1alpha1.Bucket {
	var myBucketClaim *v1alpha1.BucketClaim

	err := retry(ctx, attempts, sleep, func() error {
		var err error
		myBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
		if err != nil {
			return err
		}

		if myBucketClaim.Status.BucketName == "" {
			return fmt.Errorf("BucketName is empty")
		}
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim.Status.BucketName).NotTo(gomega.BeEmpty())

	var bucket *v1alpha1.Bucket

	err = retry(ctx, attempts, sleep, func() error {
		var err error
		bucket, err = bucketClient.ObjectstorageV1alpha1().Buckets().Get(ctx, myBucketClaim.Status.BucketName, v1.GetOptions{})
		if !bucket.Status.BucketReady {
			return fmt.Errorf("bucket %s is not ready", bucket.Name)
		}

		return err
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())

	return bucket
}

// CheckBucketStatusEmpty function for checking if Bucket status is empty.
func CheckBucketStatusEmpty(ctx context.Context, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	var myBucketClaim *v1alpha1.BucketClaim

	err := retry(ctx, attempts, sleep, func() error {
		var err error
		myBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
		return err
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim.Status.BucketName).To(gomega.BeEmpty())
}
