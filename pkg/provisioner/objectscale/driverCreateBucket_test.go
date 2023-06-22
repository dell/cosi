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
	"context"
	"testing"
	"time"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/goobjectscale/pkg/client/api/mocks"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDriverCreateBucket tests bucket creation functionality on ObjectScale platform using mock.
//
// Deprecated: this is old test suite, that is going to be removed soon.
func TestDriverCreateBucket_deprecated(t *testing.T) {
	t.Parallel()

	// Namespace (ObjectstoreID) and testID (driver ID) provided in the config file
	const (
		namespace = "namespace"
		testID    = "test.id"
	)

	testCases := []struct {
		description   string
		inputName     string
		expectedError error
		server        Server
		parameters    map[string]string
	}{ /*
			{
				description:   "valid bucket creation",
				inputName:     "bucket-valid",
				expectedError: nil,
				server: Server{
					mgmtClient: fake.NewClientSet(),
					namespace:  namespace,
					backendID:  testID,
				},
				parameters: map[string]string{
					"clientID": testID,
				},
			},
			{
				description:   "bucket already exists",
				inputName:     "bucket-valid",
				expectedError: nil,
				server: Server{
					mgmtClient: fake.NewClientSet(&model.Bucket{
						Name:      "bucket-valid",
						Namespace: namespace,
					}),
					namespace: namespace,
					backendID: testID,
				},
				parameters: map[string]string{
					"clientID": testID,
				},
			},
			{
				description:   "invalid bucket name",
				inputName:     "",
				expectedError: status.Error(codes.InvalidArgument, "empty bucket name"),
				server: Server{
					mgmtClient: fake.NewClientSet(),
					namespace:  namespace,
					backendID:  testID,
				},
				parameters: map[string]string{
					"clientID": testID,
				},
			},
			{
				description:   "cannot get existing bucket",
				inputName:     "bucket-valid",
				expectedError: status.Error(codes.Internal, "An unexpected error occurred"),
				server: Server{
					mgmtClient: fake.NewClientSet(),
					namespace:  namespace,
					backendID:  testID,
				},
				parameters: map[string]string{
					"clientID":                      testID,
					"X-TEST/Buckets/Get/force-fail": "abc",
				},
			},
			{
				description:   "cannot create bucket",
				inputName:     "FORCEFAIL-bucket-valid",
				expectedError: status.Error(codes.Internal, "failed to create bucket"),
				server: Server{
					mgmtClient: fake.NewClientSet(),
					namespace:  namespace,
					backendID:  testID,
				},
				parameters: map[string]string{
					"clientID": testID,
				},
			},
		*/}

	for _, scenario := range testCases {
		scenario := scenario

		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := scenario.server.DriverCreateBucket(ctx, &cosi.DriverCreateBucketRequest{Name: scenario.inputName, Parameters: scenario.parameters})
			assert.ErrorIs(t, err, scenario.expectedError, err)
		})
	}
}

// TestGetBucket contains table tests for (*Server).DriverCreateBucket method.
func TestDriverCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

// TestGetBucket contains table tests for (*Server).getBucket method.
func TestGetBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testValidCreateBucket": testValidCreateBucket,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testValidCreateBucket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100) // Magic value,... abra kadabra
	defer cancel()

	bucketsMock := &mocks.BucketsInterface{}
	bucketsMock.On("Create", mock.Anything, mock.Anything).Return(&model.Bucket{}, nil).Once()

	/* Boilerplate for test starts here, it should look the same for all tests. */
	mgmtClientMock := &mocks.ClientSet{}
	mgmtClientMock.On("Buckets").Return(bucketsMock).Once()

	server := Server{
		mgmtClient:    mgmtClientMock,
		namespace:     namespace,
		backendID:     testID,
		objectScaleID: objectScaleID,
		objectStoreID: objectStoreID,
	}

	err := server.createBucket(ctx, bucket)
	/* Boilerplate for test ends here, it should look the same for all tests. */

	assert.NoError(t, err)
}

// TestCreateBucket contains table tests for (*Server).createBucket method.
func TestCreateBucket(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}
