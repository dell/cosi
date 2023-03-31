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
	"regexp"
	"testing"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestExactlyOne(t *testing.T) {
	testCases := []struct {
		name      string
		nillables []interface{}
		expected  bool
	}{
		{
			name:      "exactly one non nil",
			nillables: []interface{}{nil, "x", nil},
			expected:  true,
		},
		{
			name:      "nil nillables",
			nillables: nil,
			expected:  false,
		},
		{
			name:      "empty nillables",
			nillables: []interface{}{},
			expected:  false,
		},
		{
			name:      "nil pointer",
			nillables: []interface{}{(*config.Objectscale)(nil)},
			expected:  false,
		},
		{
			name:      "all nil",
			nillables: []interface{}{nil, nil, nil},
			expected:  false,
		},
		{
			name:      "more than one not nil",
			nillables: []interface{}{nil, "x", 1},
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

var expectedOne = regexp.MustCompile("^expected exactly one OSP in configuration$")

func TestNewVirtualDriver(t *testing.T) {
	testCases := []struct {
		name         string
		config       config.Configuration
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name:   "valid config",
			config: validConfig,
		},
		{
			name:         "invalid config",
			config:       invalidConfig,
			fail:         true,
			errorMessage: expectedOne,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vd, err := NewVirtualDriver(tc.config)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
				return
			}
			assert.NoError(t, err)
			if assert.NotNil(t, vd) {
				assert.Equal(t, tc.config.Objectscale.Id, vd.ID())
			}
		})
	}
}
