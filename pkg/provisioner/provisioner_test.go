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
	"context"
	"testing"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

func TestNew(t *testing.T) {
	provSrv := New()
	if provSrv == nil {
		t.Error("server should not be nil")
	}
}

// FIXME: those are only smoke tests, no real testing is done here
func TestServer(t *testing.T) {
	provSrv := Server{}

	for scenario, fn := range map[string]func(t *testing.T, srv Server){
		"smoke/testDriverCreateBucket":       testDriverCreateBucket,
		"smoke/testDriverDeleteBucket":       testDriverDeleteBucket,
		"smoke/testDriverGrantBucketAccess":  testDriverGrantBucketAccess,
		"smoke/testDriverRevokeBucketAccess": testDriverRevokeBucketAccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, provSrv)
		})
	}
}

// FIXME: write valid test
func testDriverCreateBucket(t *testing.T, srv Server) {
	_, err := srv.DriverCreateBucket(context.TODO(), &cosi.DriverCreateBucketRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test
func testDriverDeleteBucket(t *testing.T, srv Server) {
	_, err := srv.DriverDeleteBucket(context.TODO(), &cosi.DriverDeleteBucketRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test
func testDriverGrantBucketAccess(t *testing.T, srv Server) {
	_, err := srv.DriverGrantBucketAccess(context.TODO(), &cosi.DriverGrantBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}

// FIXME: write valid test
func testDriverRevokeBucketAccess(t *testing.T, srv Server) {
	_, err := srv.DriverRevokeBucketAccess(context.TODO(), &cosi.DriverRevokeBucketAccessRequest{})
	if err == nil {
		t.Error("expected error")
	}
}
