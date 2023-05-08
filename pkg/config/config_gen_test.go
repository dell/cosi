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

package config

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	missingField  = regexp.MustCompile(`field (.+) in (.+): required`)
	invalidObject = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go value of type (.+)$`)
	invalidField  = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go struct field (.+) of type (.+)$`)
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
				"namespace":"testnamespace",
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
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
				"namespace":"testnamespace",
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
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
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing objectscale-gateway",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"namespace":"testnamespace",
				"objectstore-gateway":"gateway.objectstore.test",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing objectstore-gateway",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"namespace":"testnamespace",
				"objectscale-gateway":"gateway.objectscale.test",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing protocols",
			data: []byte(
				`{"credentials":{"username":"testuser","password":"testpassword"},
				"id":"testid",
				"namespace":"testnamespace",
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
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
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
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
