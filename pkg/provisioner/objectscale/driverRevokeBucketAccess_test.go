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
	"github.com/dell/cosi/pkg/internal/testcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamapimock "github.com/dell/cosi/pkg/internal/iamapi/mock"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
)

var _ iamiface.IAMAPI = (*iamapimock.MockIAMAPI)(nil)

func TestServerBucketAccessRevoke(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testValidAccessRevoking":              testValidAccessRevoking,
		"testNothingToChange":                  testNothingToChange,
		"testEmptyBucketIDRevoke":              testEmptyBucketIDRevoke,
		"testInvalidBucketID":                  testInvalidBucketID,
		"testEmptyAccountID":                   testEmptyAccountID,
		"testGetBucketUnexpectedError":         testGetBucketUnexpectedError,
		"testGetBucketFailToCheckUser":         testGetBucketFailToCheckUser,
		"testFailToGetAccessKeysList":          testFailToGetAccessKeysList,
		"testFailToDeleteAccessKey":            testFailToDeleteAccessKey,
		"testFailToCheckBucketPolicyExistence": testFailToCheckBucketPolicyExistence,
		"testEmptyPolicy":                      testEmptyPolicy,
		"testFailedToDeleteUser":               testFailedToDeleteUser,
		"testFailedToUpdateBucketPolicy":       testFailedToUpdateBucketPolicy,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testValidAccessRevoking tests the happy path of the (*Server).DriverRevokeBucketAccess method.
func testValidAccessRevoking(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// That's how we can mock the objectscale IAM api client
	IAMClient := iamapimock.NewMockIAMAPI(t)

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}

	IAMClient.On("ListAccessKeys", mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: accessKeyList,
	}, nil).Once()
	IAMClient.On("DeleteAccessKey", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("DeleteUser", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(testPolicy, nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	response, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

// testNothingToChange tests if no error appear when there is no resource to delete.
func testNothingToChange(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	// skip deleting access keys
	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(nil, errors.New("NoSuchEntity")).Once()

	// skip updating policy
	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrParameterNotFound).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	response, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)
	assert.ErrorIs(t, err, nil, err)
	assert.NotNil(t, response)
}

// testEmptyBucketIDRevoke tests if error handling for empty BucketID in the (*Server).DriverRevokeBucketAccess method.
func testEmptyBucketIDRevoke(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverRevokeBucketAccessRequest{
		BucketId:  "",
		AccountId: testUserName,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, req)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrInvalidBucketID.Error()))
}

func testInvalidBucketID(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverRevokeBucketAccessRequest{
		BucketId:  "bucket",
		AccountId: testUserName,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, req)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrInvalidBucketID.Error()))
}

// testEmptyAccountID tests if error handling for empty AccountID in the (*Server).DriverRevokeBucketAccess method.
func testEmptyAccountID(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	req := &cosi.DriverRevokeBucketAccessRequest{
		BucketId:  "bucket-valid",
		AccountId: "",
	}

	_, err := server.DriverRevokeBucketAccess(ctx, req)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, ErrEmpyAccountID.Error()))
}

// testGetBucketUnknownError tests if the unexpected error returned from mocked API,
// is handled correctly in the (*Server).getBucket method.
func testGetBucketUnexpectedError(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckBucketExists.Error()))
}

// testGetBucketFailToCheckUser tests if user non-existence during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testGetBucketFailToCheckUser(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(nil, ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckUserExists.Error()))
}

// testFailToGetAccessKeysList tests if failing to get access keys for a user during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testFailToGetAccessKeysList(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	// skip updating policy
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrParameterNotFound).Once()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()
	IAMClient.On("ListAccessKeys", mock.Anything).Return(nil, ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Times(3)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToListAccessKeys.Error()))
}

// testFailToDeleteAccessKey tests if failing to delete access keys for a user during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testFailToDeleteAccessKey(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	// skip updating policy
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrParameterNotFound).Once()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}
	IAMClient.On("ListAccessKeys", mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: accessKeyList,
	}, nil).Once()
	IAMClient.On("DeleteAccessKey", mock.Anything).Return(nil, ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Times(3)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToDeleteAccessKey.Error()))
}

// testFailToCheckBucketPolicyExistence tests if failing to check for policy existence during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testFailToCheckBucketPolicyExistence(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", ErrInternalException).Once()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToCheckPolicyExists.Error()))
}

// testEmptyPolicy tests if policy emptiness during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testEmptyPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}

	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrExistingPolicyIsEmpty.Error()))
}

// testFailedToDeleteUser tests if failing to delete user during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testFailedToDeleteUser(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)

	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}

	IAMClient.On("ListAccessKeys", mock.Anything).Return(&iam.ListAccessKeysOutput{
		AccessKeyMetadata: accessKeyList,
	}, nil).Once()
	IAMClient.On("DeleteAccessKey", mock.Anything).Return(nil, nil).Once()
	IAMClient.On("DeleteUser", mock.Anything).Return(nil, ErrInternalException).Once()
	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	// skip updating policy
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, ErrParameterNotFound).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToDeleteUser.Error()))
}

// testFailedToUpdateBucketPolicy tests if failing to update policy during revoking access is handled correctly
// in the (*Server).DriverRevokeAccess method.
func testFailedToUpdateBucketPolicy(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	IAMClient := iamapimock.NewMockIAMAPI(t)
	accessKeyList := make([]*iam.AccessKeyMetadata, 1)
	accessKeyList[0] = &iam.AccessKeyMetadata{
		AccessKeyId: aws.String("abc"),
		UserName:    aws.String(testUserName),
	}

	IAMClient.On("GetUser", mock.Anything).Return(&iam.GetUserOutput{
		User: &iam.User{
			UserName: aws.String(testUserName),
		},
	}, nil).Once()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{
		Name:      "valid",
		Namespace: testNamespace,
	}, nil).Once()
	bucketsMock.On("GetPolicy", mock.Anything, mock.Anything, mock.Anything).Return(testPolicy, nil).Once()
	bucketsMock.On("UpdatePolicy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock)

	server := Server{
		mgmtClient:    mgmtClientMock,
		iamClient:     IAMClient,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverRevokeBucketAccess(ctx, testBucketRevokeAccessRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, ErrFailedToUpdateBucketPolicy.Error()))
}

func TestRemove(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		inputList     []string
		inputToRemove string
		expected      []string
	}{
		{
			name:          "basic",
			inputList:     []string{"a", "b", "c"},
			inputToRemove: "a",
			expected:      []string{"b", "c"},
		},
		{
			name:          "empty input",
			inputList:     []string{},
			inputToRemove: "a",
			expected:      []string{},
		},
		{
			name:          "repeated removal",
			inputList:     []string{"a", "b", "a", "a", "c"},
			inputToRemove: "a",
			expected:      []string{"b", "c"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := remove(tc.expected, tc.inputToRemove)
			assert.Equal(t, tc.expected, output)
		})
	}
}
