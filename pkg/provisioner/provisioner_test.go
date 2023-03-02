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

package provisioner

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/fake"
	"github.com/emcecs/objectscale-management-go-sdk/pkg/client/model"
	"github.com/stretchr/testify/assert"
)

// FIXME: those are only smoke tests, no real testing is done here
func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
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

// FIXME: write valid test
func testDriverCreateBucket(t *testing.T) {
	const (
		namespace = "namespace"
		testID    = "test.id"
	)

	testCases := []struct {
		description   string
		inputName     string
		expectedError error
		server        Server
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
		},
		{
			description:   "bucket already exists",
			inputName:     "bucket-valid",
			expectedError: status.Error(codes.Internal, "Bucket already exists"),
			server: Server{
				mgmtClient: fake.NewClientSet(&model.Bucket{
					Name:      "bucket-valid",
					Namespace: namespace,
				}),
				namespace: namespace,
				backendID: testID,
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
		},
		{
			description:   "invalid bucket namespace",
			inputName:     "bucket-valid",
			expectedError: status.Error(codes.Internal, "Bucket was not sucessfully created"),
			server: Server{
				mgmtClient: fake.NewClientSet(),
				namespace:  "",
				backendID:  testID,
			},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			parameters := map[string]string{
				"clientID": testID,
			}

			_, err := scenario.server.DriverCreateBucket(context.TODO(), &cosi.DriverCreateBucketRequest{Name: scenario.inputName, Parameters: parameters})
			assert.ErrorIs(t, err, scenario.expectedError)
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
