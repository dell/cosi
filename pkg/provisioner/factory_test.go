// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package provisioner

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"

	"github.com/dell/cosi/pkg/config"
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
	testNamespace = "testnamespace"

	validConfig = config.Configuration{
		Objectscale: &config.Objectscale{
			Id:        "valid.id",
			Namespace: &testNamespace,
			Region:    aws.String("us-east-1"),
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

var expectedOne = regexp.MustCompile("^configuration is empty$")

func TestNewVirtualDriver(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]func(*testing.T){
		//		"valid config":   testValidConfig, // TODO: fix
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
	vd, err := NewVirtualDriver(validConfig)
	assert.NotNil(t, vd)
	assert.NoError(t, err)
	assert.Equal(t, validConfig.Objectscale.Id, vd.ID())
}

func testInvalidConfig(t *testing.T) {
	vd, err := NewVirtualDriver(invalidConfig)
	assert.Nil(t, vd)
	assert.Error(t, err)
	assert.Regexp(t, expectedOne, err.Error())
}
