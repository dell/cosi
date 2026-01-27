// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package config

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

var (
	missingField      = regexp.MustCompile(`field (.+) in (.+): required`)
	invalidObject     = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go value of type (.+)$`)
	invalidObjectYAML = regexp.MustCompile(`cannot unmarshal (.+) (.+) into map`)
	invalidField      = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go struct field (.+) of type (.+)$`)
)

func TestObjectscaleUnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid insecure objectscale",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"mgmt-endpoint":"https://example.com/api/s3",
				"namespace":"testnamespace",
				"emptyBucket": false,
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail: false,
		},
		{
			name:         "empty objectscale",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing credentials",
			data: []byte(
				`{"id":"testid",
				"mgmt-endpoint":"https://example.com/api/s3",
				"namespace":"testnamespace",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing id",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"namespace":"testnamespace",
				"mgmt-endpoint":"https://example.com/api/s3",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "valid missing namespace",
			data: []byte(
				`{"id":"testid",
				"credentials":{"username":"testuser","password":"testpassword"},
				"mgmt-endpoint":"https://example.com/api/s3",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         false,
			errorMessage: missingField,
		},
		{
			name: "invalid missing protocols",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"mgmt-endpoint":"https://example.com/api/s3",
				"namespace":"testnamespace",
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing tls",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"namespace":"testnamespace",
				"mgmt-endpoint":"https://example.com/api/s3",
				"protocols":{"s3":{"endpoint":"test.endpoint"}}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing mgmt-endpoint",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"namespace":"testnamespace",
				"mgmt-endpoint":"https://example.com/api/s3",
				"protocols":{"s3":{"endpoint":"test.endpoint"}}}`),
			fail:         true,
			errorMessage: missingField,
		},

		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var objectscale Objectscale

			err := objectscale.UnmarshalJSON(tc.data)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectscaleUnmarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid insecure objectscale",
			data: []byte(`---
id: testid
namespace: testnamespace
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
tls:
  insecure: true
credentials:
  username: testuser
  password: testpassword`),
			fail: false,
		},
		{
			name:         "empty objectscale",
			data:         []byte(``),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing credentials",
			data: []byte(`---
id: testid
namespace: testnamespace
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
tls:
  insecure: true`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing id",
			data: []byte(`---
namespace: testnamespace
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
tls:
  insecure: true
credentials:
  username: testuser
  password: testpassword`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "valid missing namespace",
			data: []byte(`---
id: testid
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
tls:
  insecure: true
credentials:
  username: testuser
  password: testpassword`),
			fail:         false,
			errorMessage: missingField,
		},
		{
			name: "invalid missing mgmt-endpoint",
			data: []byte(`---
id: testid
namespace: testnamespace
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
tls:
  insecure: true
credentials:
  username: testuser
  password: testpassword`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing protocols",
			data: []byte(`---
id: testid
namespace: testnamespace
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
tls:
  insecure: true
credentials:
  username: testuser
  password: testpassword`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing tls",
			data: []byte(`---
id: testid
namespace: testnamespace
mgmt-endpoint: https://example.com/api/s3
emptyBucket: false
protocols:
  s3:
    endpoint: test.endpoint
credentials:
  username: testuser
  password: testpassword`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObjectYAML,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var objectscale Objectscale
			var node yaml.Node

			err := yaml.Unmarshal(tc.data, &node)
			if err != nil {
				log.Fatalf("Error unmarshaling YAML: %v", err)
			}
			err = objectscale.UnmarshalYAML(&node)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTlsUnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid insecure",
			data: []byte(`{"insecure":true}`),
			fail: false,
		},
		{
			name:         "empty value",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var tls Tls

			err := tls.UnmarshalJSON(tc.data)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTlsUnmarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid insecure",
			data: []byte(`insecure: true`),
			fail: false,
		},
		{
			name:         "empty value",
			data:         []byte(`insecure: `),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObjectYAML,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var tls Tls
			var node yaml.Node

			err := yaml.Unmarshal(tc.data, &node)
			if err != nil {
				log.Fatalf("Error unmarshaling YAML: %v", err)
			}
			err = tls.UnmarshalYAML(&node)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid S3",
			data: []byte(`{"endpoint":"test.endpoint"}`),
			fail: false,
		},
		{
			name:         "empty endpoint",
			data:         []byte(`{"endpoint":{}}`),
			fail:         true,
			errorMessage: invalidField,
		},
		{
			name:         "missing endpoint",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "unmarshall error",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var s3 S3

			err := s3.UnmarshalJSON(tc.data)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3UnmarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid S3",
			data: []byte(`endpoint: test.endpoint`),
			fail: false,
		},
		{
			name:         "missing endpoint",
			data:         []byte(``),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "unmarshall error",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObjectYAML,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var s3 S3
			var node yaml.Node

			err := yaml.Unmarshal(tc.data, &node)
			if err != nil {
				log.Fatalf("Error unmarshaling YAML: %v", err)
			}
			err = s3.UnmarshalYAML(&node)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCredentialsUnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid credentials",
			data: []byte(`{"username":"testuser","password":"testpassword"}`),
			fail: false,
		},
		{
			name:         "missing password",
			data:         []byte(`{"username":"testuser"}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "missing username",
			data:         []byte(`{"password":"testpassword"}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var credentials Credentials

			err := credentials.UnmarshalJSON(tc.data)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCredentialsUnmarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid credentials",
			data: []byte(`username: testuser
password: testpassword`),
			fail: false,
		},
		{
			name:         "missing password",
			data:         []byte(`username: testuser`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "missing username",
			data:         []byte(`password: testpassword`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObjectYAML,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var credentials Credentials
			var node yaml.Node

			err := yaml.Unmarshal(tc.data, &node)
			if err != nil {
				log.Fatalf("Error unmarshaling YAML: %v", err)
			}
			err = credentials.UnmarshalYAML(&node)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigSchemaJsonUnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "empty config",
			data: []byte(`{}`),
			fail: false,
		},
		{
			name: "empty connections",
			data: []byte(`{"connections":[]}`),
			fail: false,
		},
		{
			name:         "invalid type",
			data:         []byte(`""`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var config ConfigSchemaJson

			err := json.Unmarshal(tc.data, &config)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
