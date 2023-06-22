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
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dell/cosi-driver/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestExactlyOne tests the exactlyOne function
// which is used to validate that only one OSP
// is defined in the configuration.
func TestExactlyOne(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		nillables []any
		expected  bool
	}{
		{
			name:      "exactly one non nil",
			nillables: []any{nil, "x", nil},
			expected:  true,
		},
		{
			name:      "nil nillables",
			nillables: nil,
			expected:  false,
		},
		{
			name:      "empty nillables",
			nillables: []any{},
			expected:  false,
		},
		{
			name:      "nil pointer",
			nillables: []any{(*config.Objectscale)(nil)},
			expected:  false,
		},
		{
			name:      "all nil",
			nillables: []any{nil, nil, nil},
			expected:  false,
		},
		{
			name:      "more than one not nil",
			nillables: []any{nil, "x", 1},
			expected:  false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := exactlyOne(tc.nillables...)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

var (
	validConfig = config.Configuration{
		Objectscale: &config.Objectscale{
			Id:                 "valid.id",
			ObjectscaleGateway: "gateway.objectscale.test",
			ObjectstoreGateway: "gateway.objectstore.test",
			ObjectscaleId:      "objectscale123",
			ObjectstoreId:      "objectstore123",
			Namespace:          "testnamespace",
			Region:             aws.String("us-east-1"),
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
		},
	}
	invalidConfig = config.Configuration{
		Objectscale: nil,
	}
)

var expectedOne = regexp.MustCompile("^expected exactly one object storage platform in configuration$")

func TestNewVirtualDriver(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(*testing.T){
		"valid config":   testValidConfig,
		"invalid config": testInvalidConfig,
	} {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			test(t)
		})
	}
}

func testValidConfig(t *testing.T) {
	vd, err := NewVirtualDriver(context.TODO(), validConfig)
	assert.NotNil(t, vd)
	assert.NoError(t, err)
	assert.Equal(t, validConfig.Objectscale.Id, vd.ID())
}

func testInvalidConfig(t *testing.T) {
	vd, err := NewVirtualDriver(context.TODO(), invalidConfig)
	assert.Nil(t, vd)
	assert.Error(t, err)
	assert.Regexp(t, expectedOne, err.Error())
}
