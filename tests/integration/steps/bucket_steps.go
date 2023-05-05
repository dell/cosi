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
	"fmt"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	errors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/container-object-storage-interface-api/apis/objectstorage/v1alpha1"
	bucketclientset "sigs.k8s.io/container-object-storage-interface-api/client/clientset/versioned"
)

// CreateBucketClaimResource Function creating a BucketClaim resource from specification.
func CreateBucketClaimResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	kubernetesBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		kubernetesBucketClaim, err = bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Create(ctx, bucketClaim, v1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	} else {
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}

	ginkgo.GinkgoWriter.Printf("Kubernetes BucketClaim: %+v\n", kubernetesBucketClaim)
}

// DeleteBucketClaimResource Function for deleting BucketClaim resource.
func DeleteBucketClaimResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	fmt.Printf("bucketClaim.Namespace: %v", bucketClaim.Namespace)
	err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Delete(ctx, bucketClaim.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketClaimStatus Function for checking BucketClaim status.
func CheckBucketClaimStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim, status bool) {
	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim).NotTo(gomega.BeNil())
	gomega.Expect(myBucketClaim.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketStatus Function for checking Bucket status.
func CheckBucketStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket, status bool) {
	gomega.Expect(bucket.Status.BucketReady).To(gomega.Equal(status))
}

// CheckBucketID Function for checking bucketID.
func CheckBucketID(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucket *v1alpha1.Bucket) {
	gomega.Expect(bucket.Status.BucketID).NotTo(gomega.Or(gomega.BeEmpty(), gomega.BeNil()))
}

// CreateBucketClassResource Function for creating BucketClass resource.
func CreateBucketClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Get(ctx, bucketClass.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		bucketclass, err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Create(ctx, bucketClass, v1.CreateOptions{})
		fmt.Printf("Error: %v", err)
		fmt.Printf("Bucketclass: %v", bucketclass)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	} else {
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}
}

// DeleteBucketClassResource Function for deleting BucketClass resource.
func DeleteBucketClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClass *v1alpha1.BucketClass) {
	err := bucketClient.ObjectstorageV1alpha1().BucketClasses().Delete(ctx, bucketClass.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucketAccessClassResource Function for creating BucketAccessClass resource.
func CreateBucketAccessClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Get(ctx, bucketAccessClass.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Create(ctx, bucketAccessClass, v1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	} else {
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}
}

// DeleteBucketAccessClassResource Function for deleting BucketAccessClass resource.
func DeleteBucketAccessClassResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccessClass *v1alpha1.BucketAccessClass) {
	err := bucketClient.ObjectstorageV1alpha1().BucketAccessClasses().Delete(ctx, bucketAccessClass.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucketAccessResource Function for creating BucketAccess resource.
func CreateBucketAccessResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	_, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Create(ctx, bucketAccess, v1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	} else {
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}
}

// DeleteBucketAccessResource Function for deleting BucketAccess resource.
func DeleteBucketAccessResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess) {
	err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Delete(ctx, bucketAccess.Name, v1.DeleteOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketAccessStatus Function for checking BucketAccess status.
func CheckBucketAccessStatus(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, status bool) {
	myBucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(myBucketAccess.Status.AccessGranted).To(gomega.Equal(status))
}

// CheckBucketAccessAccountID Function for checking BucketAccess accountID.
func CheckBucketAccessAccountID(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketAccess *v1alpha1.BucketAccess, accountID string) {
	myBucketAccess, err := bucketClient.ObjectstorageV1alpha1().BucketAccesses(bucketAccess.Namespace).Get(ctx, bucketAccess.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketAccess).NotTo(gomega.BeNil())
	gomega.Expect(myBucketAccess.Status.AccountID).To(gomega.Equal(accountID))
}

// CheckBucketResource Function for getting Bucket resource.
func GetBucketResource(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) *v1alpha1.Bucket {

	var myBucketClaim *v1alpha1.BucketClaim

	err := retry(ctx, 5, 2, func() error { // nolint:gomnd
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

	err = retry(ctx, 5, 2, func() error { // nolint:gomnd
		var err error
		bucket, err = bucketClient.ObjectstorageV1alpha1().Buckets().Get(ctx, myBucketClaim.Status.BucketName, v1.GetOptions{})
		return err
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bucket).NotTo(gomega.BeNil())

	ginkgo.GinkgoWriter.Printf("Kubernetes Bucket: %+v\n", bucket)

	return bucket
}

func CheckBucketStatusEmpty(ctx ginkgo.SpecContext, bucketClient *bucketclientset.Clientset, bucketClaim *v1alpha1.BucketClaim) {
	// Wait for creations of bucket in cluster
	time.Sleep(2 * time.Second) // nolint:gomnd

	myBucketClaim, err := bucketClient.ObjectstorageV1alpha1().BucketClaims(bucketClaim.Namespace).Get(ctx, bucketClaim.Name, v1.GetOptions{})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(myBucketClaim.Status.BucketName).To(gomega.BeEmpty())
}

func retry(ctx ginkgo.SpecContext, attempts int, sleep time.Duration, f func() error) error {
	ticker := time.NewTicker(sleep)
	retries := 0
	for {
		select {
		case <-ticker.C:
			err := f()
			if err == nil {
				return nil
			}

			retries++
			if retries > attempts {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
