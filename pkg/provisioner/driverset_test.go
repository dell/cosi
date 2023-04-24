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

package provisioner

import (
	"sync"
	"testing"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	"github.com/dell/cosi-driver/pkg/provisioner/virtualdriver/fake"
	"github.com/stretchr/testify/assert"
)

func TestDriversetAdd(t *testing.T) {
	t.Parallel()

	driverset := &Driverset{}
	driverset.drivers.Store("driver0", &fake.Driver{FakeID: "driver0"})

	testCases := []struct {
		name         string
		driverset    *Driverset
		driver       driver.Driver
		want         *Driverset
		wantErrorMsg string
	}{
		{
			name:      "no duplicate",
			driverset: &Driverset{},
			driver:    &fake.Driver{FakeID: "driver0"},
			want:      driverset,
		},
		{
			name:         "duplicate",
			driverset:    driverset,
			driver:       &fake.Driver{FakeID: "driver0"},
			want:         driverset,
			wantErrorMsg: "failed to load new driver to driverset sync.Map",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.driverset.Add(tc.driver)
			if err != nil {
				assert.ErrorContains(t, err, tc.wantErrorMsg)
			}
			compareSyncMaps(t, &tc.want.drivers, &tc.driverset.drivers)
		})
	}
}

func compareSyncMaps(t *testing.T, want, got *sync.Map) {
	t.Helper()

	wantNormal := make(map[string]driver.Driver)

	want.Range(func(key, value interface{}) bool {
		wantNormal[key.(string)] = value.(driver.Driver)
		return true
	})

	gotNormal := make(map[string]driver.Driver)

	got.Range(func(key, value interface{}) bool {
		gotNormal[key.(string)] = value.(driver.Driver)
		return true
	})

	assert.Equal(t, wantNormal, gotNormal)
}

func TestInvalidType(t *testing.T) {
	t.Parallel()

	driverset := &Driverset{}
	driverset.drivers.Store("invalid", "value")

	_, err := driverset.Get("invalid")
	assert.Error(t, err)
}

func TestDriversetGet(t *testing.T) {
	t.Parallel()

	driverset := &Driverset{}
	driverset.drivers.Store("driver0", &fake.Driver{FakeID: "driver0"})

	testCases := []struct {
		name         string
		driverset    *Driverset
		id           string
		want         driver.Driver
		wantErrorMsg string
	}{
		{
			name:      "driver configured",
			driverset: driverset,
			id:        "driver0",
			want:      &fake.Driver{FakeID: "driver0"},
		},
		{
			name:         "driver not configured",
			driverset:    driverset,
			id:           "driver1",
			want:         nil,
			wantErrorMsg: "failed to get driver from driverset",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.driverset.Get(tc.id)
			if err != nil {
				assert.ErrorContains(t, err, tc.wantErrorMsg)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrDriverDuplicate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "error prints correctly",
			id:   "driverID",
			want: "driver for 'driverID' already exists",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ErrDriverDuplicate{tc.id}
			assert.Equal(t, err.Error(), tc.want)
		})
	}
}

func TestErrNotConfigured(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "error prints correctly",
			id:   "driverID",
			want: "platform identified by 'driverID' was not configured",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ErrNotConfigured{tc.id}
			assert.Equal(t, err.Error(), tc.want)
		})
	}
}
