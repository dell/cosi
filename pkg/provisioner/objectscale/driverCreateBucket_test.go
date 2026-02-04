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

// TestServerDriverCreateBucket contains table tests for (*Server).DriverCreateBucket method.
func TestServerDriverCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		// happy path
		"BucketCreated": testDriverCreateBucketBucketCreated,
		"BucketExists":  testDriverCreateBucketBucketExists,
		// testing errors
		"CheckBucketFailed":    testDriverCreateBucketCheckBucketFailed,
		"BucketCreationFailed": testDriverCreateBucketBucketCreationFailed,
		"InvalidQuotaLimit":    testDriverCreateBucketInvalidQuotaLimit,
		"VPool List Fails":     testDriverCreateBucketVPoolFails,
		"VPool Does Not Exist": testDriverCreateBucketVPoolDoesNotExist,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

const (
	testBucketName = "test_bucket"
	testNamespace  = "namespace"
	testID         = "test.id"
)

var (
	testBucket = &model.Bucket{
		Namespace: testNamespace,
		Name:      testBucketName,
	}

	testBucketCreationWithVPoolRequest = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
		Parameters: map[string]string{
			"protocol":         "S3",
			"namespace":        testNamespace,
			"replicationGroup": "rg1",
		},
	}

	testBucketCreationRequestInvalidQuotaLimit = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
		Parameters: map[string]string{
			"protocol":   "S3",
			"namespace":  testNamespace,
			"quotaLimit": "string",
		},
	}

	testBucketCreationRequest = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
		Parameters: map[string]string{
			"protocol":  "S3",
			"namespace": testNamespace,
		},
	}

	testBucketCreationRequestEmptyNamespace = &cosi.DriverCreateBucketRequest{
		Name: testBucketName,
		Parameters: map[string]string{
			"protocol":  "S3",
			"namespace": "",
		},
	}
)

func testDriverCreateBucketInvalidQuotaLimit(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	mgmtClientMock := mocks.NewClientSet(t)

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	res, err := server.DriverCreateBucket(ctx, testBucketCreationRequestInvalidQuotaLimit)

	assert.Error(t, err)
	assert.Nil(t, res)
}

// testDriverCreateBucketBucketCreated tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the bucket does not exist on the backend.
func testDriverCreateBucketBucketCreated(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(testBucket, nil).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.ErrParameterNotFound).Once()

	vpool := mocks.NewVPoolServiceInterface(t)
	vpool.On("List", mock.Anything).Return([]model.DataServiceVPool{
		{
			Name: "rg1",
		},
	}, nil).Once()
	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()
	mgmtClientMock.On("VPools").Return(vpool)

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	expectedBucketID := strings.Join([]string{server.backendID, testBucket.Name}, "-")

	res, err := server.DriverCreateBucket(ctx, testBucketCreationWithVPoolRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.BucketId, expectedBucketID)
}

func testDriverCreateBucketVPoolDoesNotExist(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	vpool := mocks.NewVPoolServiceInterface(t)
	vpool.On("List", mock.Anything).Return([]model.DataServiceVPool{
		{
			Name: "rg2",
		},
	}, nil)
	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("VPools").Return(vpool).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	res, err := server.DriverCreateBucket(ctx, testBucketCreationWithVPoolRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func testDriverCreateBucketVPoolFails(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	vpool := mocks.NewVPoolServiceInterface(t)
	vpool.On("List", mock.Anything).Return(nil, errors.New("error"))

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("VPools").Return(vpool).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	res, err := server.DriverCreateBucket(ctx, testBucketCreationWithVPoolRequest)

	assert.Error(t, err)
	assert.Nil(t, res)
}

// testDriverCreateBucketBucketExists tests the happy path of the (*Server).DriverCreateBucket method.
// It assumes that the bucket already exists on the backend.
func testDriverCreateBucketBucketExists(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(testBucket, nil).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	expectedBucketID := strings.Join([]string{server.backendID, testBucket.Name}, "-")

	res, err := server.DriverCreateBucket(ctx, testBucketCreationRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.BucketId, expectedBucketID)
}

// testDriverCreateBucketCheckBucketFailed tests if error during checking bucket existence is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucketCheckBucketFailed(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "custom error")).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	_, err := server.DriverCreateBucket(ctx, testBucketCreationRequest)

	assert.ErrorIs(t, status.Error(codes.Internal, "error finding bucket"), err)
}

// testDriverCreateBucketBucketCreationFailed tests if error during creation of bucket is handled correctly
// in the (*Server).DriverCreateBucket method.
func testDriverCreateBucketBucketCreationFailed(t *testing.T) {
	ctx, cancel := testcontext.New(t)
	defer cancel()

	bucketsMock := mocks.NewBucketServiceInterface(t)
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "error")).Once()
	bucketsMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.ErrParameterNotFound).Once()

	mgmtClientMock := mocks.NewClientSet(t)
	mgmtClientMock.On("Buckets").Return(bucketsMock).Twice()

	server := Server{
		mgmtClient: mgmtClientMock,
		namespace:  testNamespace,
		backendID:  testID,
	}

	_, err := server.DriverCreateBucket(ctx, testBucketCreationRequest)

	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed to create bucket"))
}
