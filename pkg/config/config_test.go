// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package config

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	missingFile      = regexp.MustCompile(`^unable to read config file`)
	invalidExtension = regexp.MustCompile(`^invalid file extension, should be .json, .yaml or .yml$`)

	validJSON = `{
    "connections": [
        {
            "objectscale": {
                "credentials": {
                    "username": "testuser",
                    "password": "testpassword"
                },
                "id": "testid",
                "mgmt-endpoint":"https://example.com/api/s3",
                "namespace": "testnamespace",
                "emptyBucket": false,
                "protocols": {
                    "s3": {
                        "endpoint": "test.endpoint"
                    }
                },
                "tls": {
                    "insecure": true
                }
            }
        }
    ]
}`

	invalidJSON = `{
    "connections": [
        {
            "objectscale": {
                "credentials": {
                    "username": "testuser"
                },
                "id": "testid",
                "namespace": "testnamespace",
                "protocols": {
                    "s3": {
                        "endpoint": "test.endpoint"
                    }
                },
                "tls": {
                    "insecure": true
                }
            }
        }
    ]
}`

	validYAML = `connections:
- objectscale:
    credentials:
      username: testuser
      password: testpassword
    id: testid
    namespace: testnamespace
    mgmt-endpoint: https://example.com/api/s3
    emptyBucket: false
    protocols:
      s3:
        endpoint: test.endpoint
    tls:
      insecure: true`

	invalidYAML = `connections:
- objectscale:
    credentials:
      username: testuser
    id: testid
    namespace: test-namespace
    emptyBucket: false
    protocols:
      s3:
        endpoint: test.endpoint
    tls:
      insecure: true`
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		file         testFile
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "valid JSON",
			file: testFile{
				name:    "valid.json",
				content: validJSON,
			},
			fail: false,
		},
		{
			name: "invalid JSON",
			file: testFile{
				name:    "invalid.json",
				content: invalidJSON,
			},
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "valid YAML",
			file: testFile{
				name:    "valid.yaml",
				content: validYAML,
			},
			fail: false,
		},
		{
			name: "valid YML",
			file: testFile{
				name:    "valid.yml",
				content: validYAML,
			},
			fail: false,
		},
		{
			name: "invalid YAML",
			file: testFile{
				name:    "invalid.yaml",
				content: invalidYAML,
			},
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "invalid YML",
			file: testFile{
				name:    "invalid.yml",
				content: invalidYAML,
			},
			fail:         true,
			errorMessage: missingField,
		},
		{
			name: "missing JSON file",
			file: testFile{
				skip: true,
				name: "missing.json",
			},
			fail:         true,
			errorMessage: missingFile,
		},
		{
			name: "missing YAML file",
			file: testFile{
				skip: true,
				name: "missing.yaml",
			},
			fail:         true,
			errorMessage: missingFile,
		},
		{
			name: "missing YML file",
			file: testFile{
				skip: true,
				name: "missing.yml",
			},
			fail:         true,
			errorMessage: missingFile,
		},
		{
			name: "invalid file extension",
			file: testFile{
				name: "invalid.txt",
			},
			fail:         true,
			errorMessage: invalidExtension,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testFile, err := tc.file.Write()
			defer os.RemoveAll(path.Dir(testFile))

			if err != nil {
				// unexpected error, should panic
				panic(err)
			}

			config, err := New(testFile)
			if tc.fail {
				if assert.Errorf(t, err, "%+#v", config) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err, testFile)
			}
		})
	}
}

type testFile struct {
	name    string
	content string
	skip    bool
}

func (tf *testFile) Write() (string, error) {
	s, err := os.MkdirTemp("", "test-*")
	if err != nil {
		return "", err
	}

	file := path.Join(s, tf.name)

	// if file should not be written, skip
	if tf.skip == true {
		return file, nil
	}

	return file, os.WriteFile(file, []byte(tf.content), 0o600)
}
