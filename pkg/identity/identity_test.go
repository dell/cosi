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
	"io"
	"os"
	"testing"

	"github.com/dell/cosi/pkg/internal/testcontext"
	log "github.com/sirupsen/logrus"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	provisioner = "test"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

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
	srv := New(provisioner)

	ctx, cancel := testcontext.New(t)
	defer cancel()

	res, err := srv.DriverGetInfo(ctx, &cosi.DriverGetInfoRequest{})
	if err != nil {
		t.Errorf("got unexpected error: %s", err.Error())
	}

	if res.Name != provisioner {
		t.Errorf("got: '%s', wanted '%s'", res.Name, provisioner)
	}
}

func testMissingProvisionerName(t *testing.T) {
	srv := &Server{}

	ctx, cancel := testcontext.New(t)
	defer cancel()

	_, err := srv.DriverGetInfo(ctx, &cosi.DriverGetInfoRequest{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
