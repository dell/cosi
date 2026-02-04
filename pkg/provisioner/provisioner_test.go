// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package provisioner

import (
	"fmt"
	"os"
	"testing"

	cosi "sigs.k8s.io/container-object-storage-interface/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dell/cosi/pkg/internal/testcontext"
	"github.com/dell/cosi/pkg/provisioner/virtualdriver/fake"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestNew tests the initialization of provisioner server.
func TestNew(t *testing.T) {
	t.Parallel()

	fakeDriverset := &Driverset{}

	err := fakeDriverset.Add(&fake.Driver{FakeID: "fake"})
	if err != nil {
		t.Fatalf("failed to create fakedriverset: %v", err)
	}

	testServer := New(fakeDriverset)
	assert.NotNil(t, testServer)
	assert.NotNil(t, testServer.driverset)
}

// TestServer starts a server for running tests of the multi-backend provisioner.
func TestServer(t *testing.T) {
	t.Parallel()

	fakeDriverset := &Driverset{}
	err := fakeDriverset.Add(&fake.Driver{FakeID: "fake"})
	assert.Nil(t, err)

	fakeServer := Server{
		driverset: fakeDriverset,
	}

	for name, test := range map[string]func(*testing.T, Server){
		"DriverCreateBucket":       testServerDriverCreateBucket,
		"DriverDeleteBucket":       testServerDriverDeleteBucket,
		"DriverGrantBucketAccess":  testServerDriverGrantBucketAccess,
		"DriverRevokeBucketAccess": testServerDriverRevokeBucketAccess,
	} {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test(t, fakeServer)
		})
	}
}

// testServerDriverCreateBucket tests passing the DriverCreateBucketRequest to the proper driver from the driverset.
func testServerDriverCreateBucket(t *testing.T, fakeServer Server) {
	testCases := []struct {
		server        Server
		description   string
		req           *cosi.DriverCreateBucketRequest
		expectedError error
	}{
		{
			server:      fakeServer,
			description: "bucket creation successful",
			req: &cosi.DriverCreateBucketRequest{
				Name: "fake",
				Parameters: map[string]string{
					fake.KeyDriverID: "fake",
				},
			},
			expectedError: nil,
		},
		{
			server:      fakeServer,
			description: "bucket creation force fail",
			req: &cosi.DriverCreateBucketRequest{
				Name: "fake",
				Parameters: map[string]string{
					fake.KeyDriverID: "fake",
					fake.ForceFail:   "true",
				},
			},
			expectedError: status.Error(codes.Internal, "an unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket creation invalid backend ID",
			req: &cosi.DriverCreateBucketRequest{
				Name: "invalid",
				Parameters: map[string]string{
					fake.KeyDriverID: "invalid",
				},
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := testcontext.New(t)
			defer cancel()

			_, err := tc.server.DriverCreateBucket(ctx, &cosi.DriverCreateBucketRequest{
				Name:       tc.req.Name,
				Parameters: tc.req.Parameters,
			})
			assert.ErrorIs(t, err, tc.expectedError, err)
		})
	}
}

// testServerDriverDeleteBucket tests passing the DriverDeleteBucketRequest to the proper driver from the driverset.
func testServerDriverDeleteBucket(t *testing.T, fakeServer Server) {
	testCases := []struct {
		server        Server
		description   string
		req           *cosi.DriverDeleteBucketRequest
		expectedError error
	}{
		{
			server:      fakeServer,
			description: "bucket deletion successful",
			req: &cosi.DriverDeleteBucketRequest{
				BucketId: "fake-bucket",
			},
			expectedError: nil,
		},
		{
			server:      fakeServer,
			description: "bucket deletion force fail",
			req: &cosi.DriverDeleteBucketRequest{
				BucketId: fmt.Sprintf("fake-%s", fake.ForceFail),
			},
			expectedError: status.Error(codes.Internal, "an unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket deletion invalid backend ID",
			req: &cosi.DriverDeleteBucketRequest{
				BucketId: "invalid",
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := testcontext.New(t)
			defer cancel()

			_, err := tc.server.DriverDeleteBucket(ctx, &cosi.DriverDeleteBucketRequest{
				BucketId: tc.req.BucketId,
			})
			assert.ErrorIs(t, err, tc.expectedError, err)
		})
	}
}

// testServerDriverGrantBucketAccess tests passing the DriverGrantBucketAccessRequest
// to the proper driver from the driverset.
func testServerDriverGrantBucketAccess(t *testing.T, fakeServer Server) {
	testCases := []struct {
		server        Server
		description   string
		req           *cosi.DriverGrantBucketAccessRequest
		expectedError error
	}{
		{
			server:      fakeServer,
			description: "bucket access granting successful",
			req: &cosi.DriverGrantBucketAccessRequest{
				BucketId: "valid-bucket",
				Parameters: map[string]string{
					fake.KeyDriverID: "fake",
				},
			},
			expectedError: nil,
		},
		{
			server:      fakeServer,
			description: "bucket access granting failed",
			req: &cosi.DriverGrantBucketAccessRequest{
				BucketId: "valid-bucket",
				Parameters: map[string]string{
					fake.KeyDriverID: "fake",
					fake.ForceFail:   "true",
				},
			},
			expectedError: status.Error(codes.Internal, "an unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket access granting invalid backend ID",
			req: &cosi.DriverGrantBucketAccessRequest{
				Name: "invalid",
				Parameters: map[string]string{
					fake.KeyDriverID: "invalid",
				},
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := testcontext.New(t)
			defer cancel()

			_, err := tc.server.DriverGrantBucketAccess(ctx, &cosi.DriverGrantBucketAccessRequest{
				BucketId:   tc.req.BucketId,
				Parameters: tc.req.Parameters,
			})
			assert.ErrorIs(t, err, tc.expectedError, err)
		})
	}
}

// testServerDriverRevokeBucketAccess tests passing the DriverRevokeBucketAccessRequest
// to the proper driver from the driverset.
func testServerDriverRevokeBucketAccess(t *testing.T, fakeServer Server) {
	testCases := []struct {
		server        Server
		description   string
		req           *cosi.DriverRevokeBucketAccessRequest
		expectedError error
	}{
		{
			server:      fakeServer,
			description: "bucket access revoking successful",
			req: &cosi.DriverRevokeBucketAccessRequest{
				BucketId: "fake-bucket",
			},
			expectedError: nil,
		},
		{
			server:      fakeServer,
			description: "bucket access revoking force fail",
			req: &cosi.DriverRevokeBucketAccessRequest{
				BucketId: fmt.Sprintf("fake-%s", fake.ForceFail),
			},
			expectedError: status.Error(codes.Internal, "an unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket access revoking invalid backend ID",
			req: &cosi.DriverRevokeBucketAccessRequest{
				BucketId: "invalid",
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := testcontext.New(t)
			defer cancel()

			_, err := tc.server.DriverRevokeBucketAccess(ctx, &cosi.DriverRevokeBucketAccessRequest{
				BucketId: tc.req.BucketId,
			})
			assert.ErrorIs(t, err, tc.expectedError, err)
		})
	}
}
