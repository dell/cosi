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
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dell/goobjectscale/pkg/client/fake"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/internal/testcontext"
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
		"testNew":                  testDriverNew,
		"testID":                   testDriverID,
		"testParsePolicyStatement": testParsePolicyStatement,
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
			_, cancel := testcontext.New(t)
			defer cancel()

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

func testParsePolicyStatement(t *testing.T) {
	testCases := []struct {
		description          string
		inputStatements      []UpdateBucketPolicyStatement
		awsBucketResourceARN string
		awsPrincipalString   string
		expectedOutput       []UpdateBucketPolicyStatement
	}{
		{
			description: "valid policy statement parsing",
			inputStatements: []UpdateBucketPolicyStatement{
				{
					Resource: []string{"happyAwsBucketResourceARN"},
					SID:      "GetObject_permission",
					Effect:   allowEffect,
					Principal: Principal{
						AWS: []string{"happyAwsPrincipalString"},
					},
					Action: []string{"*"},
				},
			},
			awsBucketResourceARN: "happyAwsBucketResourceARN",
			awsPrincipalString:   "happyAwsPrincipalString",
			expectedOutput: []UpdateBucketPolicyStatement{
				{
					Resource: []string{"happyAwsBucketResourceARN"},
					SID:      "GetObject_permission",
					Effect:   allowEffect,
					Principal: Principal{
						AWS: []string{"happyAwsPrincipalString"},
					},
					Action: []string{"*"},
				},
			},
		},
		{
			description: "policy needed update parsing",
			inputStatements: []UpdateBucketPolicyStatement{
				{
					Resource: nil,
					SID:      "GetObject_permission",
					Effect:   "",
					Principal: Principal{
						AWS: []string{"urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118"},
					},
					Action: []string{"s3:GetObjectVersion"},
				},
			},
			awsBucketResourceARN: "happyAwsBucketResourceARN",
			awsPrincipalString:   "happyAwsPrincipalString",
			expectedOutput: []UpdateBucketPolicyStatement{
				{
					Resource: []string{"happyAwsBucketResourceARN"},
					SID:      "GetObject_permission",
					Effect:   allowEffect,
					Principal: Principal{
						AWS: []string{"urn:osc:iam::osai07c2ae318ae9d6f2:user/iam_user20230523061025118", "happyAwsPrincipalString"},
					},
					Action: []string{"s3:GetObjectVersion", "*"},
				},
			},
		},
	}

	for _, scenario := range testCases {
		t.Run(scenario.description, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			updatedPolicy := parsePolicyStatement(ctx, scenario.inputStatements, scenario.awsBucketResourceARN, scenario.awsPrincipalString)
			assert.Equalf(t, scenario.expectedOutput, updatedPolicy, "not equal")
		})
	}
}

// TestGetBucketName tests BucketID splitting.
func TestGetBucketName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic",
			input:    "first-second",
			expected: "second",
		},
		{
			name:     "extra_dashes",
			input:    "first-second-third",
			expected: "second-third",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output, err := GetBucketName(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, output)
		})
	}
}

func FuzzGetBucketName(f *testing.F) {
	for _, seed := range []string{
		"driverid-bucketname",
		"driver.id-bucket.name",
		".driver.id-bucket.name",
		".driver.id-bucket-name",
		"--",
		"-.",
	} {
		f.Add(seed) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, in string) {
		out, err := GetBucketName(in)

		// Check if the input string has one hyphen, but not ending with a hyphen
		hasHyphen := strings.Count(in, "-") == 1 && !strings.HasSuffix(in, "-")
		// Check if the input string has more than one hyphen
		hasMultipleHyphens := strings.Count(in, "-") > 1

		// If the input has one hyphen and doesn't end with a hyphen, or has multiple hyphens
		if hasHyphen || hasMultipleHyphens {
			// The error should be nil
			assert.NoErrorf(t, err, "Input was '%s'", in)
			// The output should not be empty
			assert.NotEmpty(t, out, "Input was '%s'", in)
		} else {
			// Otherwise, there should be an error
			assert.Errorf(t, err, "Input was '%s'", in)
		}
	})
}
