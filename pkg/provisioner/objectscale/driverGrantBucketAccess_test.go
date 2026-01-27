// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/dell/cosi/pkg/internal/testcontext"
	omocks "github.com/dell/cosi/pkg/provisioner/objectscale/mocks"
	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

var testBucketGrantAccessRequest = &cosi.DriverGrantBucketAccessRequest{
	BucketId: strings.Join([]string{testID, testBucketName}, "-"),
	Name:     "bucket-access-id",
}

func TestServerDriverGrantAccess(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"GrantAccess":           testDriverGrantBucketAccess,
		"GrantAccessUserExists": testDriverGrantBucketAccessUserExists,
		// testing errors
		"UnableToGetIAMClient":                    testDriverGrantBucketAccessUnableToGetIAMClient,
		"ErrorCheckingBucketExistence":            testDriverGrantBucketAccessErrorCheckingBucketExistence,
		"BucketAccessBucketDoesNotExist":          testDriverGrantBucketAccessBucketDoesNotExist,
		"GrantAccessErrorGettingIAMUser":          testDriverGrantBucketAccessErrorGettingIAMUser,
		"GrantAccessErrorGettingBucketPolicy":     testDriverGrantBucketAccessErrorGettingBucketPolicy,
		"GrantAccessErrorInvalidJSONPolicy":       testDriverGrantBucketAccessErrorInvalidJSONPolicy,
		"GrantBucketAccessErrorUpdatingPolicy":    testDriverGrantBucketAccessErrorUpdatingPolicy,
		"GrantBucketAccessErrorCreatingAccessKey": testDriverGrantBucketAccessErrorCreatingAccessKey,
		"GrantBucketAccessErrorCreatingUser":      testDriverGrantBucketAccessErrorCreatingUser,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testDriverGrantBucketAccess(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(&iam.CreateUserOutput{
		User: &types.User{
			UserName: aws.String("user"),
		},
	}, nil).Once()
	iamMock.On("CreateAccessKey", mock.Anything, mock.Anything).Return(&iam.CreateAccessKeyOutput{
		AccessKey: &types.AccessKey{
			AccessKeyId:     aws.String("key"),
			SecretAccessKey: aws.String("secret"),
		},
	}, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "namespace-user-bucket-access-id", res.AccountId)
}

func testDriverGrantBucketAccessErrorGettingIAMUser(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.MalformedPolicyDocumentException{}).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessUserExists(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{User: &types.User{
		UserName: aws.String("user"),
	}}, nil).Once()
	iamMock.On("CreateAccessKey", mock.Anything, mock.Anything).Return(&iam.CreateAccessKeyOutput{
		AccessKey: &types.AccessKey{
			AccessKeyId:     aws.String("key"),
			SecretAccessKey: aws.String("secret"),
		},
	}, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "namespace-user-bucket-access-id", res.AccountId)
}

func testDriverGrantBucketAccessErrorCheckingBucketExistence(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error finding bucket")).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessBucketDoesNotExist(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.ErrParameterNotFound).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessUnableToGetIAMClient(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	mgmtClientMock := mocks.NewClientSet(t)

	val := func(context.Context) (IAM, error) {
		return nil, errors.New("iam client error")
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessErrorGettingBucketPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("error getting bucket policy")).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(&iam.CreateUserOutput{
		User: &types.User{
			UserName: aws.String("user"),
		},
	}, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessErrorInvalidJSONPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("invalid-json", nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(&iam.CreateUserOutput{
		User: &types.User{
			UserName: aws.String("user"),
		},
	}, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessErrorUpdatingPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error updating policy")).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(&iam.CreateUserOutput{
		User: &types.User{
			UserName: aws.String("user"),
		},
	}, nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessErrorCreatingAccessKey(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(&iam.CreateUserOutput{
		User: &types.User{
			UserName: aws.String("user"),
		},
	}, nil).Once()
	iamMock.On("CreateAccessKey", mock.Anything, mock.Anything).Return(nil, errors.New("error creating access key")).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverGrantBucketAccessErrorCreatingUser(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	iamMock := omocks.NewIAM(t)
	iamMock.On("GetUser", mock.Anything, mock.Anything).Return(nil, &types.NoSuchEntityException{}).Once()
	iamMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("error creating user"), nil).Once()

	val := func(context.Context) (IAM, error) {
		return iamMock, nil
	}
	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
		iamClient:  val,
	}

	res, err := server.DriverGrantBucketAccess(ctx, testBucketGrantAccessRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}
