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

package pkg

import (
	"context"
	"testing"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	provisioner = "test"
)

// FIXME: those are only smoke tests, no real testing is done here
func TestDriverGetInfo(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"smoke/testValidServer":            testValidServer,
		"smoke/testMissingProvisionerName": testMissingProvisionerName,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// FIXME: write valid test
func testValidServer(t *testing.T) {
	srv := &IdentityServer{
		provisioner: provisioner,
	}

	res, err := srv.DriverGetInfo(context.TODO(), &cosi.DriverGetInfoRequest{})
	if err != nil {
		t.Errorf("got unexpected error: %s", err.Error())
	}

	if res.Name != provisioner {
		t.Errorf("got: '%s', wanted '%s'", res.Name, provisioner)
	}
}

// FIXME: write valid test
func testMissingProvisionerName(t *testing.T) {
	srv := &IdentityServer{}

	_, err := srv.DriverGetInfo(context.TODO(), &cosi.DriverGetInfoRequest{})
	if err == nil {
		t.Errorf("got unexpected error: %s", err.Error())
	}
}
