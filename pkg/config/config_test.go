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
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDir = "test"
)

var (
	dir string
)

var (
	missingFile      = regexp.MustCompile(`^unable to read config file: open (.*): no such file or directory$`)
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
                "objectscale-gateway": "gateway.objectscale.test",
                "objectstore-gateway": "gateway.objectstore.test",
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
                "objectscale-gateway": "gateway.objectscale.test",
                "objectstore-gateway": "gateway.objectstore.test",
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
    objectscale-gateway: gateway.objectscale.test
    objectstore-gateway: gateway.objectstore.test
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
    objectscale-gateway: gateway.objectscale.test
    objectstore-gateway: gateway.objectstore.test
    protocols:
      s3:
        endpoint: test.endpoint
    tls:
      insecure: true`
)

func TestNew(t *testing.T) {
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

	// create test dir
	var err error
	dir, err = os.MkdirTemp("", testDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.file.Write()
			if err != nil {
				// unexpected error, should panic
				panic(err)
			}

			testfile := path.Join(dir, tc.file.name)

			x, err := New(testfile)
			if tc.fail {
				if assert.Errorf(t, err, "%+#v", x) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type testFile struct {
	name    string
	content string
	skip    bool
}

func (tf *testFile) Write() error {
	// if file should not be written, skip
	if tf.skip == true {
		return nil
	}

	return os.WriteFile(path.Join(dir, tf.name), []byte(tf.content), 0644)
}
