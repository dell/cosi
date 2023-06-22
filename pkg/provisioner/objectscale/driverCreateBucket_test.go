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

	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

// TestGetBucket contains table tests for (*Server).DriverCreateBucket method.
func TestDriverCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"BucketCreated": testDriverCreateBucketBucketCreated,
		"BucketExists":  testDriverCreateBucketBucketExists,
		// testing errors
		"EmptyBucketName":      testDriverCreateBucketEmptyBucketName,
		"CheckBucketFailed":    testDriverCreateBucketCheckBucketFailed,
		"BucketCreationFailed": testDriverCreateBucketBucketCreationFailed,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testDriverCreateBucketBucketCreated tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the driver does not exist on the backend.
func testDriverCreateBucketBucketCreated(t *testing.T) {
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

// testDriverCreateBucketBucketExists tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the driver already exists on the backend.
func testDriverCreateBucketBucketExists(t *testing.T) {
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

// testDriverCreateBucketEmptyBucketName tests if missing bucket name is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucketEmptyBucketName(t *testing.T) {
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

// testDriverCreateBucketCheckBucketFailed tests if error during checking bucket existence is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucketCheckBucketFailed(t *testing.T) {
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

// testDriverCreateBucketBucketCreationFailed tests if error during creation of bucket is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucketBucketCreationFailed(t *testing.T) {
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
		"Valid": testGetBucketValid,
		// testing errors
		"NoBucket":     testGetBucketNoBucket,
		"UnknownError": testGetBucketUnknownError,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testGetBucketValid tests the happy path of the (*Server).getBucket method.
func testGetBucketValid(t *testing.T) {
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

// testGetBucketNoBucket tests if the error indicating that no bucket was found returned from the mocked API,
// is handled correctly in the (*Server).getBucket method.
func testGetBucketNoBucket(t *testing.T) {
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

// testGetBucketUnknownError tests if the unexpected error returned from mocked API,
// is handled correctly in the (*Server).getBucket method.
func testGetBucketUnknownError(t *testing.T) {
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
		"Valid": testCreateBucketValid,
		// testing errors
		"Failed": testCreateBucketFailed,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testCreateBucketValid tests the happy path of the (*Server).createBucket method.
func testCreateBucketValid(t *testing.T) {
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

// testCreateBucketValid tests if the error returned from the mocked API is handled correctly
// in the (*Server).createBucket method.
func testCreateBucketFailed(t *testing.T) {
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
