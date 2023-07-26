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
	"testing"

	"github.com/dell/cosi-driver/pkg/internal/testcontext"
	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestServerDriverDeleteBucket contains table tests for (*Server).DriverDeleteBucket method.
func TestServerDriverDeleteBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"BucketDeleted":       testDriverDeleteBucketBucketDeleted,
		"BucketDoesNotExists": testDriverDeleteBucketBucketDoesNotExists,
		// // testing errors
		"InvalidBucketID":      testDriverDeleteBucketInvalidBucketID,
		"BucketDeletionFailed": testDriverDeleteBucketBucketDeletionFailed,
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

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	res, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.NoError(t, err)
	require.NotNil(t, res)
}

// testDriverDeleteBucketBucketDoesNotExists tests the happy path of the (*Server).DriverDeleteBucket method.
// It assumes that bucket does not exist on the backend.
func testDriverDeleteBucketBucketDoesNotExists(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ErrParameterNotFound).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	res, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.NoError(t, err)
	require.NotNil(t, res)
}

// testDriverDeleteBucketInvalidBucketID tests if missing bucket ID is handled correctly.
// in the (*Server).DriverDeleteBucket method.
func testDriverDeleteBucketInvalidBucketID(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	mgmtClientMock := &mocks.ClientSet{}

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequestEmptyBucketID)

	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "empty bucketID"))
}

// testDriverDeleteBucketBucketDeleted tests if error during deletion of bucket is handled correctly
// in the (*Server).DriverDeleteBucket method.
func testDriverDeleteBucketBucketDeletionFailed(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ErrInternalException).Once()

	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     testNamespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	_, err := server.DriverDeleteBucket(ctx, testBucketDeletionRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "bucket was not successfully deleted"))
}
