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

package config

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	missingField  = regexp.MustCompile(`field (.+) in (.+): required`)
	invalidEnum   = regexp.MustCompile(`^invalid value \(expected one of (.+)\): (.+)$`)
	invalidObject = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go value of type (.+)$`)
	invalidField  = regexp.MustCompile(`^json: cannot unmarshal (.+) into Go struct field (.+) of type (.+)$`)
)

func TestTlsMinVersionUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid TLS 1.3",
			data: []byte(`"TLS-1.3"`),
			fail: false,
		},
		{
			name: "valid TLS 1.2",
			data: []byte(`"TLS-1.2"`),
			fail: false,
		},
		{
			name: "valid TLS 1.1",
			data: []byte(`"TLS-1.1"`),
			fail: false,
		},
		{
			name: "valid TLS 1.0",
			data: []byte(`"TLS-1.0"`),
			fail: false,
		},
		{
			name:         "invalid TLS",
			data:         []byte(`"unknown"`),
			fail:         true,
			errorMessage: invalidEnum,
		},
		{
			name:         "empty TLS",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var version TlsMinVersion

			err := version.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectscaleUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid insecure objectscale",
			data: []byte(
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
				"id":"testid",
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
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
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
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
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
				"id":"testid",
				"objectstore-gateway":"gateway.objectstore.test",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing objectstore-gateway",
			data: []byte(
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
				"id":"testid",
				"objectscale-gateway":"gateway.objectscale.test",
				"protocols":{"s3":{"endpoint":"test.endpoint"}},
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing protocols",
			data: []byte(
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
				"id":"testid",
				"objectscale-gateway":"gateway.objectscale.test",
				"objectstore-gateway":"gateway.objectstore.test",
				"tls":{"insecure":true}}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid missing tls",
			data: []byte(
				`{"credentials":{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="},
				"id":"testid",
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
		t.Run(tc.name, func(t *testing.T) {
			var objectscale Objectscale

			err := objectscale.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTlsUnmarshalJSON(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			var tls Tls

			err := tls.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigSchemaJsonLogLevelUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid level(fatal)",
			data: []byte(`"fatal"`),
			fail: false,
		},
		{
			name: "valid level(error)",
			data: []byte(`"error"`),
			fail: false,
		},
		{
			name: "valid level(warning)",
			data: []byte(`"warning"`),
			fail: false,
		},
		{
			name: "valid level(info)",
			data: []byte(`"info"`),
			fail: false,
		},
		{
			name: "valid level(debug)",
			data: []byte(`"debug"`),
			fail: false,
		},
		{
			name: "valid level(trace)",
			data: []byte(`"trace"`),
			fail: false,
		},
		{
			name:         "invalid value",
			data:         []byte(`"unknown"`),
			fail:         true,
			errorMessage: invalidEnum,
		},
		{
			name:         "invalid type",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var logLevel ConfigSchemaJsonLogLevel

			err := logLevel.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTlsClientAuthUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid tls.NoClientCert",
			data: []byte(`"tls.NoClientCert"`),
			fail: false,
		},
		{
			name: "valid tls.RequestClientCert",
			data: []byte(`"tls.RequestClientCert"`),
			fail: false,
		},
		{
			name: "valid tls.RequireAnyClientCert",
			data: []byte(`"tls.RequireAnyClientCert"`),
			fail: false,
		},
		{
			name: "valid tls.VerifyClientCertIfGiven",
			data: []byte(`"tls.VerifyClientCertIfGiven"`),
			fail: false,
		},
		{
			name: "valid tls.RequireAndVerifyClientCert",
			data: []byte(`"tls.RequireAndVerifyClientCert"`),
			fail: false,
		},
		{
			name:         "invalid value",
			data:         []byte(`"unknown"`),
			fail:         true,
			errorMessage: invalidEnum,
		},
		{
			name:         "invalid type",
			data:         []byte(`{}`),
			fail:         true,
			errorMessage: invalidObject,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var clientAuth TlsClientAuth

			err := clientAuth.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3UnmarshalJSON(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			var s3 S3

			err := s3.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCredentialsUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid credentials",
			data: []byte(`{"username":"dGVzdHVzZXIK","password":"dGVzdHBhc3N3b3JkCg=="}`),
			fail: false,
		},
		{
			name:         "missing password",
			data:         []byte(`{"username":"dGVzdHVzZXIK"}`),
			fail:         true,
			errorMessage: missingField,
		},
		{
			name:         "missing username",
			data:         []byte(`{"password":"dGVzdHBhc3N3b3JkCg=="}`),
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
		t.Run(tc.name, func(t *testing.T) {
			var credentials Credentials

			err := credentials.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigSchemaJsonUnmarshalJSON(t *testing.T) {
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
			name: "filled log-level",
			data: []byte(`{"log-level":"info"}`),
			fail: false,
		},
		{
			name: "filled cosi-endpoint",
			data: []byte(`{"cosi-endpoint":"unix:///var/lib/cosi/cosi.sock"}`),
			fail: false,
		},
		{
			name: "filled cosi-endpoint, log-level",
			data: []byte(`{"cosi-endpoint":"unix:///var/lib/cosi/cosi.sock","log-level":"info"}`),
			fail: false,
		},
		{
			name: "filled cosi-endpoint, log-level, present empty connections",
			data: []byte(`{"connections":[],"cosi-endpoint":"unix:///var/lib/cosi/cosi.sock","log-level":"info"}`),
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
		t.Run(tc.name, func(t *testing.T) {
			var config ConfigSchemaJson

			err := config.UnmarshalJSON(tc.data)
			if tc.fail {
				assert.Error(t, err)
				assert.Regexp(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
