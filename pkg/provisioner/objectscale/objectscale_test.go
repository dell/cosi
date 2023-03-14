//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/fake"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
)

// FIXME: those are only smoke tests, no real testing is done here
func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"smoke/testNew":                      testDriverNew,
		"smoke/testID":                       testDriverID,
		"smoke/testDriverCreateBucket":       testDriverCreateBucket,
		"smoke/testDriverDeleteBucket":       testDriverDeleteBucket,
		"smoke/testDriverGrantBucketAccess":  testDriverGrantBucketAccess,
		"smoke/testDriverRevokeBucketAccess": testDriverRevokeBucketAccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testDriverNew tests server initialization
func testDriverNew(t *testing.T) {
	driver := Server{
		mgmtClient: fake.NewClientSet(),
		backendID:  "id",
		namespace:  "namespace",
	}
	assert.Equal(t, "id", driver.backendID)
	assert.Equal(t, "namespace", driver.namespace)
	assert.NotNil(t, driver.mgmtClient)
}

// testDriverID tests extending COSI interface by adding driver ID
func testDriverID(t *testing.T) {
	driver := Server{
		mgmtClient: fake.NewClientSet(),
		backendID:  "id",
		namespace:  "namespace",
	}
	assert.Equal(t, "id", driver.ID())
}

// testDriverCreateBucket tests bucket creation functionality on ObjectScale platform
func testDriverCreateBucket(t *testing.T) {
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
	}{
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
			expectedError: status.Error(codes.AlreadyExists, "Bucket already exists"),
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
			expectedError: status.Error(codes.InvalidArgument, "Empty bucket name"),
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
			expectedError: status.Error(codes.Internal, "Bucket was not sucessfully created"),
			server: Server{
				mgmtClient: fake.NewClientSet(),
				namespace:  namespace,
				backendID:  testID,
			},
			parameters: map[string]string{
				"clientID": testID,
			},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			_, err := scenario.server.DriverCreateBucket(context.TODO(), &cosi.DriverCreateBucketRequest{Name: scenario.inputName, Parameters: scenario.parameters})
			assert.ErrorIs(t, err, scenario.expectedError, err)
		})
	}
}

// FIXME: write valid test
func testDriverDeleteBucket(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverDeleteBucket(context.TODO(), &cosi.DriverDeleteBucketRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test
func testDriverGrantBucketAccess(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverGrantBucketAccess(context.TODO(), &cosi.DriverGrantBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test
func testDriverRevokeBucketAccess(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverRevokeBucketAccess(context.TODO(), &cosi.DriverRevokeBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}
