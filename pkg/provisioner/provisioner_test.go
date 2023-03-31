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

package provisioner

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	"github.com/dell/cosi-driver/pkg/provisioner/virtualdriver/fake"
)

// TestNew tests the initialization of provisioner server.
func TestNew(t *testing.T) {
	t.Parallel()

	fakeDriverset := &Driverset{drivers: map[string]virtualdriver.Driver{}}

	err := fakeDriverset.Add(&fake.Driver{FakeID: "fake"})
	if err != nil {
		log.Fatalf("Failed to create fakedriverset: %v", err)
	}

	testServer := New(fakeDriverset)

	testServer := New(fakeDriverset)
	assert.NotNil(t, testServer)
	assert.NotNil(t, testServer.driverset)
}

// TestServer starts a server for running tests of the multi-backend provisioner.
func TestServer(t *testing.T) {
	t.Parallel()

	fakeDriverset := &Driverset{drivers: map[string]virtualdriver.Driver{}}
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
			expectedError: status.Error(codes.Internal, "An unexpected error occurred"),
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
			expectedError: status.Error(codes.InvalidArgument, "DriverCreateBucket: Invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			_, err := tc.server.DriverCreateBucket(context.TODO(), &cosi.DriverCreateBucketRequest{
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
			expectedError: status.Error(codes.Internal, "An unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket deletion invalid backend ID",
			req: &cosi.DriverDeleteBucketRequest{
				BucketId: "invalid",
			},
			expectedError: status.Error(codes.InvalidArgument, "DriverDeleteBucket: Invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			_, err := tc.server.DriverDeleteBucket(context.TODO(), &cosi.DriverDeleteBucketRequest{
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
			expectedError: status.Error(codes.Internal, "An unexpected error occurred"),
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
			expectedError: status.Error(codes.InvalidArgument, "DriverGrantBucketAccess: Invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			_, err := tc.server.DriverGrantBucketAccess(context.TODO(), &cosi.DriverGrantBucketAccessRequest{
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
			expectedError: status.Error(codes.Internal, "An unexpected error occurred"),
		},
		{
			server:      fakeServer,
			description: "bucket access revoking invalid backend ID",
			req: &cosi.DriverRevokeBucketAccessRequest{
				BucketId: "invalid",
			},
			expectedError: status.Error(codes.InvalidArgument, "DriverRevokeBucketAccess: Invalid backend ID"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			_, err := tc.server.DriverRevokeBucketAccess(context.TODO(), &cosi.DriverRevokeBucketAccessRequest{
				BucketId: tc.req.BucketId,
			})
			assert.ErrorIs(t, err, tc.expectedError, err)
		})
	}
}
