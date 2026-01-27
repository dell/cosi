// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package identity

import (
	"os"
	"testing"

	"github.com/dell/cosi/pkg/internal/testcontext"
	cosi "sigs.k8s.io/container-object-storage-interface/proto"
)

const (
	provisioner = "test"
)

func TestMain(m *testing.M) {
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
