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
	"errors"
	"fmt"
	"reflect"

	"github.com/dell/cosi/pkg/provisioner/objectscale"
	"github.com/dell/cosi/pkg/provisioner/policy"
	"github.com/dell/goobjectscale/pkg/client/api"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	smithy "github.com/aws/smithy-go"

	"sigs.k8s.io/container-object-storage-interface/client/apis/objectstorage/v1alpha1"

	gomega "github.com/onsi/gomega"
)

// CheckObjectScaleInstallation Ensure that ObjectScale platform is installed on the cluster.
func CheckObjectScaleInstallation(ctx context.Context, client api.ClientSet, namespace string) {
	parameters := map[string]string{}
	parameters["namespace"] = namespace

	_, err := client.Buckets().List(ctx, parameters)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketResourceInObjectStore Function checking if Bucket resource is in objectstore.
func CheckBucketResourceInObjectStore(ctx context.Context, client api.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	id, err := objectscale.GetBucketNameFromID(bucket.Status.BucketID)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	parameters := map[string]string{}
	parameters["namespace"] = namespace

	checkAttempt := 0
	err = retry(ctx, attempts, sleep, func() error {
		checkAttempt++

		objectScaleBucket, err := client.Buckets().Get(ctx, id, parameters)
		if err != nil {
			fmt.Printf("[CheckBucketResourceInObjectStore] attempt %d/%d: error fetching bucket %s: %v\n", checkAttempt, attempts, id, err)
			return err
		}

		if objectScaleBucket == nil {
			fmt.Printf("[CheckBucketResourceInObjectStore] attempt %d/%d: bucket %s not found in ObjectStore\n", checkAttempt, attempts, id)
			return fmt.Errorf("bucket %s not found in ObjectStore", id)
		}

		fmt.Printf("[CheckBucketResourceInObjectStore] attempt %d/%d: bucket %s found in ObjectStore\n", checkAttempt, attempts, id)
		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckBucketDeletionInObjectStore Function for checking Bucket deletion in ObjectStore.
func CheckBucketDeletionInObjectStore(ctx context.Context, client api.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	id, err := objectscale.GetBucketNameFromID(bucket.Status.BucketID)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	parameters := map[string]string{}
	parameters["namespace"] = namespace

	deleteAttempt := 0
	err = retry(ctx, attempts, sleep, func() error {
		var err error
		var bucket *model.Bucket

		deleteAttempt++

		bucket, err = client.Buckets().Get(ctx, id, parameters)
		// if error is ErrParameterNotFound, it means the bucket was deleted from ObjectScale
		if err != nil && errors.Is(err, model.ErrParameterNotFound) {
			fmt.Printf("[CheckBucketDeletionInObjectStore] attempt %d/%d: bucket %s confirmed deleted\n", deleteAttempt, attempts, id)
			return nil
		}

		if err != nil {
			fmt.Printf("[CheckBucketDeletionInObjectStore] attempt %d/%d: error fetching bucket %s: %v\n", deleteAttempt, attempts, id, err)
			return err
		}

		if bucket != nil {
			fmt.Printf("[CheckBucketDeletionInObjectStore] attempt %d/%d: bucket %s still exists\n", deleteAttempt, attempts, bucket.Name)
			return fmt.Errorf("bucket %s still exists", bucket.Name)
		}

		return nil
	})

	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CheckUser checks if user exists in ObjectScale.
func CheckUser(ctx context.Context, iamClient *iam.Client, accountID string) {
	userOut, err := iamClient.GetUser(ctx, &iam.GetUserInput{UserName: &accountID})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(userOut.User).NotTo(gomega.BeNil())
	gomega.Expect(accountID).To(gomega.Equal(*(userOut.User.UserName)))
}

func CheckPolicy(ctx context.Context, mgmtClient api.ClientSet, expectedPolicyDocument policy.Document, myBucket *v1alpha1.Bucket, namespace string) {
	var actualPolicyDocument policy.Document

	ErrComparisonFailed := errors.New("comparison failed")

	policyAttempt := 0
	// This also needs to be retried, as we are not sure, if the policy was already updated.
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		policyAttempt++

		param := make(map[string]string)
		param["namespace"] = namespace

		actualPolicy, err := mgmtClient.Buckets().GetPolicy(ctx, myBucket.Name, param)
		if err != nil {
			fmt.Printf("[CheckPolicy] attempt %d/%d: error fetching policy for bucket %s: %v\n", policyAttempt, attempts, myBucket.Name, err)
			return err
		}

		actualPolicyDocument, err = policy.NewFromJSON(actualPolicy)
		if err != nil {
			fmt.Printf("[CheckPolicy] attempt %d/%d: error parsing policy for bucket %s: %v\n", policyAttempt, attempts, myBucket.Name, err)
			return err
		}

		if !reflect.DeepEqual(actualPolicyDocument, expectedPolicyDocument) {
			fmt.Printf("[CheckPolicy] attempt %d/%d: policy for bucket %s does not match expected\n", policyAttempt, attempts, myBucket.Name)
			return ErrComparisonFailed
		}

		fmt.Printf("[CheckPolicy] attempt %d/%d: policy for bucket %s matches expected\n", policyAttempt, attempts, myBucket.Name)
		return nil
	})

	// If the error is ErrComparisonFailed, I want full gomega match on objects, so I get pretty output.
	gomega.Expect(err).To(gomega.Or(gomega.BeEquivalentTo(ErrComparisonFailed), gomega.Not(gomega.HaveOccurred())))
	gomega.Expect(expectedPolicyDocument).To(gomega.BeEquivalentTo(actualPolicyDocument))
}

func CheckEmptyPolicy(ctx context.Context, mgmtClient api.ClientSet, myBucket *v1alpha1.Bucket, namespace string) {
	ErrComparisonFailed := errors.New("comparison failed")

	emptyPolicyAttempt := 0
	// This also needs to be retried, as we are not sure, if the policy was already updated.
	err := retry(ctx, attempts, sleep, func() error {
		var err error

		emptyPolicyAttempt++

		param := make(map[string]string)
		param["namespace"] = namespace

		actualPolicy, err := mgmtClient.Buckets().GetPolicy(ctx, myBucket.Name, param)
		if err != nil {
			fmt.Printf("[CheckEmptyPolicy] attempt %d/%d: error fetching policy for bucket %s: %v\n", emptyPolicyAttempt, attempts, myBucket.Name, err)
			return err
		}

		if actualPolicy != "" {
			fmt.Printf("[CheckEmptyPolicy] attempt %d/%d: policy for bucket %s is not empty yet\n", emptyPolicyAttempt, attempts, myBucket.Name)
			return errors.New("policy is not empty")
		}

		fmt.Printf("[CheckEmptyPolicy] attempt %d/%d: policy for bucket %s is empty as expected\n", emptyPolicyAttempt, attempts, myBucket.Name)
		return nil
	})

	// If the error is ErrComparisonFailed, I want full gomega match on objects, so I get pretty output.
	gomega.Expect(err).To(gomega.Or(gomega.BeEquivalentTo(ErrComparisonFailed), gomega.Not(gomega.HaveOccurred())))
}

// CheckUserDeleted checks if user does not exist in ObjectScale.
func CheckUserDeleted(ctx context.Context, iamClient *iam.Client, user string, _ string) {
	_, err := iamClient.GetUser(ctx, &iam.GetUserInput{UserName: &user})
	gomega.Expect(err).To(gomega.HaveOccurred())

	var apiError smithy.APIError
	matched := errors.As(err, &apiError)

	gomega.Expect(matched).To(gomega.BeTrue())
	switch apiError.(type) {
	case *types.NoSuchEntityException:
		// expected
	default:
		gomega.Panic()
	}

	// gomega.Expect(ok).To(gomega.BeTrue())
	// gomega.Expect(types.NoSuchEntityException.ErrorCodeOverride).To(gomega.Equal(apiError.ErrorCode()))
}

// CheckBucketNotInObjectStore Function for checking if bucket is not in objectstore.
func CheckBucketNotInObjectStore(ctx context.Context, client api.ClientSet, bucketClaim *v1alpha1.BucketClaim, namespace string) {
	parameters := map[string]string{}
	parameters["namespace"] = namespace

	_, err := client.Buckets().Get(ctx, bucketClaim.Status.BucketName, parameters)
	isNotFound := errors.Is(err, model.ErrParameterNotFound)
	gomega.Expect(isNotFound).To(gomega.BeTrue())
}

// DeleteBucket Function for deleting existing from ObjectScale (useful if BucketClaim deletionPolicy is set to "retain").
func DeleteBucket(ctx context.Context, client api.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	parameters := map[string]string{}
	parameters["namespace"] = namespace

	err := client.Buckets().Delete(ctx, bucket.Name, parameters)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}

// CreateBucket Function for creating bucket on Objectscale.
func CreateBucket(ctx context.Context, client api.ClientSet, namespace string, bucket *v1alpha1.Bucket) {
	model := model.ObjectBucketParam{
		Name:      bucket.Name,
		Namespace: namespace,
	}
	nbucket, err := client.Buckets().Create(ctx, &model)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(nbucket).ToNot(gomega.BeNil())
}
