// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package objectscale

import (
	"errors"
	"strings"
	"testing"

	"github.com/dell/cosi/pkg/internal/testcontext"
	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

var testBucketDeletionRequest = &cosi.DriverDeleteBucketRequest{
	BucketId: strings.Join([]string{testID, testBucketName}, "-"),
}

// TestServerDriverDeleteBucket contains table tests for (*Server).DriverDeleteBucket method.
func TestServerDriverDeleteBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"BucketDeleted":      testDriverDeleteBucketBucketDeleted,
		"BucketDoesNotExist": testDriverDeleteBucketBucketDoesNotExist,
		// testing errors
		"BucketDeletionFailed": testDriverDeleteBucketBucketDeletionFailed,
		"GetBucketFailed":      testDriverDeleteGetBucketFailed,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// testDriverDeleteBucketBucketDeleted tests the happy path of the (*Server).DriverDeleteBucket method.
// It assumes that bucket exists on the backend.
func testDriverDeleteBucketBucketDeleted(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	res, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// testDriverDeleteBucketBucketDeletionFailed tests if error during deletion of bucket is handled correctly
// in the (*Server).DriverDeleteBucket method.
func testDriverDeleteBucketBucketDeletionFailed(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("custom")).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	_, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed deleting bucket"))
}

func testDriverDeleteGetBucketFailed(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error getting bucket")).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	_, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed checking if bucket exists"))
}

func testDriverDeleteBucketBucketDoesNotExist(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.ErrParameterNotFound).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	res, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
