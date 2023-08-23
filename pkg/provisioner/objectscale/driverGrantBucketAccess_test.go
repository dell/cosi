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

package objectscale

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamapimock "github.com/dell/cosi/pkg/internal/iamapi/mock"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi/pkg/internal/testcontext"
	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
)

var _ iamiface.IAMAPI = (*iamapimock.MockIAMAPI)(nil)

func TestServerBucketAccessGrant(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testValidAccessGranting":       testValidAccessGranting,
		"testInvalidAccessKeyCreation":  testInvalidAccessKeyCreation,
		"testInvalidUserCreation":       testInvalidUserCreation,
		"testInvalidUserRetrieval":      testInvalidUserRetrieval,
		"testInvalidBucketPolicyUpdate": testInvalidBucketPolicyUpdate,
		"testEmptyBucketID":             testEmptyBucketID,
		"testEmptyName":                 testEmptyName,
		"testInvalidAuthenticationType": testInvalidAuthenticationType,
		"testIAMNotImplemented":         testIAMNotImplemented,
		"testFailToGetBucket":           testFailToGetBucket,
		"testBucketNotFound":            testBucketNotFound,
		"testValidButUserAlreadyExists": testValidButUserAlreadyExists,
		"testFailToGetExistingPolicy":   testFailToGetExistingPolicy,
		"testInvalidPolicyJSON":         testInvalidPolicyJSON,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testValidAccessGranting(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(&iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: aws.String("acc"), SecretAccessKey: aws.String("sec")}}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

func testInvalidAccessKeyCreation(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(nil, errors.New("failed to create access key")).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCreateAccessKey.Error()), err)
	assert.Nil(t, response)
}

func testInvalidUserCreation(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		nil, errors.New(ErrFailedToCreateUser.Error())).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCreateUser.Error()), err)
	assert.Nil(t, response)
}

func testInvalidUserRetrieval(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, errors.New("failed to retrieve user")).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckUserExists.Error()), err)
	assert.Nil(t, response)
}

func testInvalidBucketPolicyUpdate(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ErrInternalException).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToUpdatePolicy.Error()), err)
	assert.Nil(t, response)
}

func testEmptyBucketID(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	mgmtClientMock := &mocks.ClientSet{}

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrInvalidBucketID.Error()), err)
	assert.Nil(t, response)
}

func testEmptyName(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	mgmtClientMock := &mocks.ClientSet{}

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrEmptyBucketAccessName.Error()), err)
	assert.Nil(t, response)
}

func testInvalidAuthenticationType(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	mgmtClientMock := &mocks.ClientSet{}

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_UnknownAuthenticationType,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrInvalidAuthenticationType.Error()), err)
	assert.Nil(t, response)
}

func testIAMNotImplemented(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	mgmtClientMock := &mocks.ClientSet{}

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_IAM,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Unimplemented, ErrAuthenticationTypeNotImplemented.Error()), err)
	assert.Nil(t, response)
}

func testFailToGetBucket(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrInternalException).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckBucketExists.Error()), err)
	assert.Nil(t, response)
}

func testBucketNotFound(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrParameterNotFound).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-invalid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.NotFound, ErrBucketNotFound.Error()), err)
	assert.Nil(t, response)
}

func testValidButUserAlreadyExists(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String("namespace-user-valid"),
		},
	}, nil).Once()
	IAMClient.On("CreateAccessKey", mock.Anything).Return(&iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: aws.String("acc"), SecretAccessKey: aws.String("sec")}}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

func testFailToGetExistingPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckPolicyExists.Error()), err)
	assert.Nil(t, response)
}

func testInvalidPolicyJSON(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("CreateUserWithContext", mock.Anything, mock.Anything).Return(
		&iam.CreateUserOutput{
			User: &iam.User{
				UserName: aws.String("namespace-user-valid"), // This mocked response is based on `namesapce` from server and bucketId from request
			},
		}, nil).Once()
	IAMClient.On("GetUserWithContext", mock.Anything, mock.Anything).Return(&iam.GetUserOutput{}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(testInvalidPolicy, nil).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient, // Inject mocked IAM client
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverGrantBucketAccessRequest{
		BucketId:           "bucket-valid",
		Name:               "bucket-access-valid",
		AuthenticationType: cosi.AuthenticationType_Key,
	}

	response, err := server.DriverGrantBucketAccess(ctx, req)
	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToDecodePolicy.Error()), err)
	assert.Nil(t, response)
}
