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
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/dell/goobjectscale/pkg/client/model"
)

type expected int

const (
	ok      expected = iota
	warning expected = iota
	fail    expected = iota
)

var (
	invalidBase64 = `💀`

	validConfig = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
	}

	invalidConfigWithHyphens = &config.Objectscale{
		Id:                 "id-with-hyphens",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
	}

	invalidConfigEmptyID = &config.Objectscale{
		Id:                 "",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: true,
		},
	}

	invalidConfigTLS = &config.Objectscale{
		Id:                 "valid.id",
		ObjectscaleGateway: "gateway.objectscale.test",
		ObjectstoreGateway: "gateway.objectstore.test",
		Credentials: config.Credentials{
			Username: "testuser",
			Password: "testpassword",
		},
		Protocols: config.Protocols{
			S3: &config.S3{
				Endpoint: "s3.objectstore.test",
			},
		},
		Tls: config.Tls{
			Insecure: false,
			RootCas:  &invalidBase64,
		},
	}
)

// regex for error messages.
var (
	emptyID             = regexp.MustCompile(`^empty id$`)
	transportInitFailed = regexp.MustCompile(`^initialization of transport failed:`)
)

// testDriverNew tests server initialization.
func testDriverNew(t *testing.T) {
	testCases := []struct {
		name         string
		config       *config.Objectscale
		result       expected
		errorMessage *regexp.Regexp
	}{
		{
			name:   "valid config",
			config: validConfig,
			result: ok,
		},
		{
			name:   "invalid config with hyphens",
			config: invalidConfigWithHyphens,
			result: warning,
		},
		{
			name:         "invalid config empty id",
			config:       invalidConfigEmptyID,
			result:       fail,
			errorMessage: emptyID,
		},
		{
			name:         "invalid config TLS error",
			config:       invalidConfigTLS,
			result:       fail,
			errorMessage: transportInitFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			driver, err := New(tc.config)
			switch tc.result {
			case ok:
				assert.NoError(t, err)
				if assert.NotNil(t, driver) {
					assert.Equal(t, tc.config.Id, driver.ID())
				}

			case warning:
				assert.NoError(t, err)
				if assert.NotNil(t, driver) {
					assert.Equal(t, strings.ReplaceAll(tc.config.Id, "-", "_"), driver.ID())
				}

			case fail:
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			}
		})
	}
}

// testDriverID tests extending COSI interface by adding driver ID.
func testDriverID(t *testing.T) {
	driver := Server{
		mgmtClient: fake.NewClientSet(),
		backendID:  "id",
		namespace:  "namespace",
	}
	assert.Equal(t, "id", driver.ID())
}

// testDriverCreateBucket tests bucket creation functionality on ObjectScale platform.
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
			expectedError: status.Error(codes.Internal, "Bucket was not sucessfully created"), // typo in goobjectscale
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := scenario.server.DriverCreateBucket(ctx, &cosi.DriverCreateBucketRequest{Name: scenario.inputName, Parameters: scenario.parameters})
			assert.ErrorIs(t, err, scenario.expectedError, err)
		})
	}
}

// FIXME: write valid test.
func testDriverDeleteBucket(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverDeleteBucket(context.TODO(), &cosi.DriverDeleteBucketRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test.
func testDriverGrantBucketAccess(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverGrantBucketAccess(context.TODO(), &cosi.DriverGrantBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test.
func testDriverRevokeBucketAccess(t *testing.T) {
	srv := Server{}

	_, err := srv.DriverRevokeBucketAccess(context.TODO(), &cosi.DriverRevokeBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}
