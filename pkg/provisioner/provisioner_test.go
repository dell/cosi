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

package provisioner

import (
	"testing"

	"github.com/dell/cosi-driver/pkg/provisioner/virtual_driver/fake"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
	}{
		// TODO: add test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: add test body
		})
	}
}

func TestServer(t *testing.T) {
	testCases := map[string]func(*testing.T, Server){
		"DriverCreateBucket":       testServerDriverCreateBucket,
		"DriverDeleteBucket":       testServerDriverDeleteBucket,
		"DriverGrantBucketAccess":  testServerDriverGrantBucketAccess,
		"DriverRevokeBucketAccess": testServerDriverRevokeBucketAccess,
	}

	driverset := &Driverset{}
	driverset.Add(&fake.Driver{})
	server := Server{
		driverset: driverset,
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			test(t, server)
		})
	}
}

func testServerDriverCreateBucket(t *testing.T, server Server) {
	testCases := []struct {
		name string
	}{
		// TODO: add test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: add test body
		})
	}
}

func testServerDriverDeleteBucket(t *testing.T, server Server) {
	testCases := []struct {
		name string
	}{
		// TODO: add test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: add test body
		})
	}
}

func testServerDriverGrantBucketAccess(t *testing.T, server Server) {
	testCases := []struct {
		name string
	}{
		// TODO: add test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: add test body
		})
	}
}

func testServerDriverRevokeBucketAccess(t *testing.T, server Server) {
	testCases := []struct {
		name string
	}{
		// TODO: add test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: add test body
		})
	}
}
