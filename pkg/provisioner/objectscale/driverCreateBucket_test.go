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
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestGetBucket contains table tests for (*Server).DriverCreateBucket method.
func TestDriverCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"testDriverCreateBucket_BucketCreated": testDriverCreateBucket_BucketCreated,
		"testDriverCreateBucket_BucketExists":  testDriverCreateBucket_BucketExists,
		// testing errors
		"testDriverCreateBucket_EmptyBucketName":      testDriverCreateBucket_EmptyBucketName,
		"testDriverCreateBucket_CheckBucketFailed":    testDriverCreateBucket_CheckBucketFailed,
		"testDriverCreateBucket_BucketCreationFailed": testDriverCreateBucket_BucketCreationFailed,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testDriverCreateBucket_BucketCreated tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the driver does not exist on the backend.
func testDriverCreateBucket_BucketCreated(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(testBucket, nil).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeParameterNotFound}).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	res, err := server.DriverCreateBucket(ctx, testRequest)

	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.BucketId, strings.Join([]string{server.backendID, testBucket.Name}, "-"))
}

// testDriverCreateBucket_BucketExists tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the driver already exists on the backend.
func testDriverCreateBucket_BucketExists(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(testBucket, nil).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	res, err := server.DriverCreateBucket(ctx, testRequest)

	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.BucketId, strings.Join([]string{server.backendID, testBucket.Name}, "-"))
}

// testDriverCreateBucket_EmptyBucketName tests if missing bucket name is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucket_EmptyBucketName(t *testing.T) {
	ctx, cancel := testContext(t)
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

	_, err := server.DriverCreateBucket(ctx, &cosi.DriverCreateBucketRequest{})

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "empty bucket name"))
}

// testDriverCreateBucket_CheckBucketFailed tests if error during checking bucket existence is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucket_CheckBucketFailed(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeInternalException}).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverCreateBucket(ctx, testRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to check if bucket exists"))
}

// testDriverCreateBucket_BucketCreationFailed tests if error during creation of bucket is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucket_BucketCreationFailed(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeInternalException}).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeParameterNotFound}).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverCreateBucket(ctx, testRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to create bucket"))
}

// TestGetBucket contains table tests for (*Server).getBucket method.
func TestGetBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"testGetBucket_Valid": testGetBucket_Valid,
		// testing errors
		"testGetBucket_NoBucket":     testGetBucket_NoBucket,
		"testGetBucket_UnknownError": testGetBucket_UnknownError,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testGetBucket_Valid tests the happy path of the (*Server).getBucket method.
func testGetBucket_Valid(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	modelBucket := &model.Bucket{Name: "valid"}
	params := make(map[string]string)

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(modelBucket, nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	bucket, err := server.getBucket(ctx, testBucketName, params)

	assert.NoError(t, err)
	assert.EqualValues(t, modelBucket, bucket)
}

// testGetBucket_NoBucket tests if the error indicating that no bucket was found returned from the mocked API,
// is handled correctly in the (*Server).getBucket method.
func testGetBucket_NoBucket(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeParameterNotFound}).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	params := make(map[string]string)

	bucket, err := server.getBucket(ctx, testBucketName, params)

	assert.Nil(t, err)
	assert.Nil(t, bucket)
}

// testGetBucket_UnknownError tests if the unexpected error returned from mocked API,
// is handled correctly in the (*Server).getBucket method.
func testGetBucket_UnknownError(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeInternalException}).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	params := make(map[string]string)

	bucket, err := server.getBucket(ctx, testBucketName, params)

	assert.ErrorIs(t, err, model.Error{Code: model.CodeInternalException})
	assert.Nil(t, bucket)
}

// TestCreateBucket contains table tests for (*Server).createBucket method.
func TestCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"testCreateBucket_Valid": testCreateBucket_Valid,
		// testing errors
		"testCreateBucket_Failed": testCreateBucket_Failed,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testCreateBucket_Valid tests the happy path of the (*Server).createBucket method.
func testCreateBucket_Valid(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	err := server.createBucket(ctx, testBucket)

	assert.NoError(t, err)
}

// testCreateBucket_Valid tests if the error returned from the mocked API is handled correctly
// in the (*Server).createBucket method.
func testCreateBucket_Failed(t *testing.T) {
	ctx, cancel := testContext(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(nil, model.Error{Code: model.CodeInternalException}).Once()

	// Generic mock for the ClientSet interface, we care only about returning Buckets from it.
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	err := server.createBucket(ctx, testBucket)

	assert.ErrorIs(t, err, model.Error{Code: model.CodeInternalException})
}
