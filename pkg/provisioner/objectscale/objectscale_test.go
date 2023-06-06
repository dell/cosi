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
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/dell/goobjectscale/pkg/client/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"

	"github.com/dell/cosi-driver/pkg/config"
)

type expected int

const (
	ok expected = iota
	warning
	fail
)

// regex for error messages.
var (
	emptyID             = regexp.MustCompile(`^empty driver id$`)
	transportInitFailed = regexp.MustCompile(`^initialization of transport failed`)
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testNew":                testDriverNew,
		"testID":                 testDriverID,
		"testDriverCreateBucket": testDriverCreateBucket,
		"testDriverDeleteBucket": testDriverDeleteBucket,
		// "testDriverGrantBucketAccess":  testDriverGrantBucketAccess,
		"testDriverRevokeBucketAccess": testDriverRevokeBucketAccess,
	} {
		fn := fn

		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

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
		{
			name:         "empty namesapce",
			config:       emptyNamespaceConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty objectstore id"),
		},
		{
			name:         "empty credentials password",
			config:       emptyPasswordConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty password"),
		},
		{
			name:         "empty credentials username",
			config:       emptyUsernameConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty username"),
		},
		{
			name:         "empty region",
			config:       emptyRegionConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty region"),
		},
		{
			name:         "region not set",
			config:       regionNotSetConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("region was not specified in config"),
		},
		{
			name:         "empty objectscale gateway",
			config:       emptyObjectscaleGatewayConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty objectscale gateway"),
		},
		{
			name:         "empty objectstore gateway",
			config:       emptyObjectstoreGatewayConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty objectstore gateway"),
		},
		{
			name:         "empty s3 endpoint",
			config:       emptyS3EndpointConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty protocol S3 endpoint"),
		},
		{
			name:         "empty objectscale id",
			config:       emptyObjectscaleIDConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty objectscaleID"),
		},
		{
			name:         "empty objectstore id",
			config:       emptyObjectstoreIDConfig,
			result:       fail,
			errorMessage: regexp.MustCompile("empty objectstoreID"),
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
			expectedError: status.Error(codes.Internal, "an unexpected error occurred"),
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
			expectedError: status.Error(codes.Internal, "bucket was not successfully created"),
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

func testDriverDeleteBucket(t *testing.T) {
	const (
		namespace = "namespace"
		testID    = "test.id"
	)

	testCases := []struct {
		description   string
		inputBucketID string
		expectedError error
		server        Server
	}{
		{
			description:   "invalid bucketID",
			inputBucketID: "",
			expectedError: status.Error(codes.InvalidArgument, "empty bucketID"),
		},
		{
			description:   "bucket does not exist",
			inputBucketID: strings.Join([]string{testID, "bucket-valid"}, "-"),
			expectedError: nil,
			server: Server{
				mgmtClient: fake.NewClientSet(),
				namespace:  namespace,
				backendID:  testID,
			},
		},
		{
			description:   "failed to delete bucket",
			inputBucketID: strings.Join([]string{testID, "bucket-invalid-FORCEFAIL"}, "-"),
			expectedError: status.Error(codes.Internal, "bucket was not successfully deleted"),
			server: Server{
				mgmtClient: fake.NewClientSet(&model.Bucket{
					Name:      "bucket-valid",
					Namespace: namespace,
				}),
				namespace:   namespace,
				backendID:   testID,
				emptyBucket: true,
			},
		},
		{
			description:   "bucket successfully deleted",
			inputBucketID: strings.Join([]string{testID, "bucket-valid"}, "-"),
			expectedError: nil,
			server: Server{
				mgmtClient: fake.NewClientSet(&model.Bucket{
					Name:      "bucket-valid",
					Namespace: namespace,
				}),
				namespace:   namespace,
				backendID:   testID,
				emptyBucket: true,
			},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := scenario.server.DriverDeleteBucket(ctx, &cosi.DriverDeleteBucketRequest{BucketId: scenario.inputBucketID})
			assert.ErrorIs(t, err, scenario.expectedError, err)
		})
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

func testParsePolicyStatement(t *testing.T) {

	testCases := []struct {
		description          string
		inputStatements      []updateBucketPolicyStatement
		awsBucketResourceARN string
		awsPrincipalString   string
		expectedOutput       []updateBucketPolicyStatement
	}{
		{
			description: "valid policy statement parsing",
			inputStatements: []updateBucketPolicyStatement{
				{
					Resource: []string{"arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:lynxbucket/*"},
					SID:      "GetObject_permission",
					Effect:   allowEffect,
					Principal: principal{
						AWS:    []string{"urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"},
						Action: []string{"s3:GetObjectVersion"},
					},
				},
			},
			awsBucketResourceARN: "happyAwsBucketResourceARN",
			awsPrincipalString:   "happyAwsPrincipalString",
			expectedOutput:       nil,
		},
		{
			description: "valid policy statement parsing",
			inputStatements: []updateBucketPolicyStatement{
				{
					Resource: []string{"arn:aws:s3:osci5b022e718aa7e0ff:osti202e682782ebcbfd:lynxbucket/*"},
					SID:      "GetObject_permission",
					Effect:   allowEffect,
					Principal: principal{
						AWS:    []string{"urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"},
						Action: []string{"s3:GetObjectVersion"},
					},
				},
			},
			awsBucketResourceARN: "",
			awsPrincipalString:   "",
			expectedError:        nil,
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			_, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			// _, err := scenario.server.DriverDeleteBucket(ctx, &cosi.DriverDeleteBucketRequest{BucketId: scenario.inputBucketID})
			assert.ErrorIs(t, nil, scenario.expectedError, nil)
		})
	}
}

// func testGeneratePolicyID(t *testing.T) {
// 	const (
// 		policyID = ""
// 	)

// 	testCases := []struct {
// 		description   string
// 		bucketName    string
// 		expectedError error
// 	}{
// 		{
// 			description:   "valid policy ID generation",
// 			bucketName:    "valid-bucket",
// 			expectedError: nil,
// 		},
// 		{
// 			description:   "invald UUID generation",
// 			bucketName:    "valid-bucket",
// 			expectedError: errors.New("failed to generate PolicyID UUID"),
// 		},
// 		{
// 			description:   "generated Policy ID was empty",
// 			bucketName:    "valid-bucket",
// 			expectedError: errors.New("generated PolicyID was empty"),
// 		},
// 	}

// 	for _, scenario := range testCases {
// 		t.Run(scenario.description, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()
// 			_, err := generatePolicyID(ctx, scenario.bucketName)
// 			assert.ErrorIs(t, err, scenario.expectedError, err)
// 		})
// 	}
// }
