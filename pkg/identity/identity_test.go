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

package identity

import (
	"context"
	"testing"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	provisioner = "test"
)

func TestDriverGetInfo(t *testing.T) {
	t.Parallel()

	for scenario, fn := range map[string]func(t *testing.T){
		"testValidServer":            testValidServer,
		"testMissingProvisionerName": testMissingProvisionerName,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			fn(t)
		})
	}
}

func testValidServer(t *testing.T) {
	t.Parallel()

	srv := New(provisioner)

	res, err := srv.DriverGetInfo(context.TODO(), &cosi.DriverGetInfoRequest{})
	if err != nil {
		t.Errorf("got unexpected error: %s", err.Error())
	}

	if res.Name != provisioner {
		t.Errorf("got: '%s', wanted '%s'", res.Name, provisioner)
	}
}

func testMissingProvisionerName(t *testing.T) {
	t.Parallel()

	srv := &Server{}

	_, err := srv.DriverGetInfo(context.TODO(), &cosi.DriverGetInfoRequest{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
