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
	"testing"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
	"github.com/dell/cosi-driver/pkg/provisioner/virtualdriver/fake"
	"github.com/stretchr/testify/assert"
)

var driverset = &Driverset{
	drivers: map[string]driver.Driver{
		"driver0": &fake.Driver{FakeID: "driver0"},
	},
}

func TestDriversetInit(t *testing.T) {
	testCases := []struct {
		name      string
		driverset *Driverset
		want      *Driverset
	}{
		{
			name:      "driverset initialised",
			driverset: &Driverset{drivers: map[string]driver.Driver{}},
			want:      &Driverset{drivers: map[string]driver.Driver{}},
		},
		{
			name:      "driverset not initialised",
			driverset: &Driverset{},
			want:      &Driverset{drivers: map[string]driver.Driver{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.driverset.init()
			assert.Equal(t, tc.want, tc.driverset)
		})
	}
}

func TestDriversetAdd(t *testing.T) {
	testCases := []struct {
		name      string
		driverset *Driverset
		driver    driver.Driver
		want      *Driverset
		wantErr   error
	}{
		{
			name:      "no duplicate",
			driverset: &Driverset{drivers: map[string]driver.Driver{}},
			driver:    &fake.Driver{FakeID: "driver0"},
			want:      driverset,
			wantErr:   nil,
		},
		{
			name:      "duplicate",
			driverset: driverset,
			driver:    &fake.Driver{FakeID: "driver0"},
			want:      driverset,
			wantErr:   ErrDriverDuplicate{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.driverset.Add(tc.driver)
			assert.IsType(t, tc.wantErr, err)
			assert.Equal(t, tc.want.drivers, tc.driverset.drivers)
		})
	}
}

func TestDriversetGet(t *testing.T) {
	testCases := []struct {
		name      string
		driverset *Driverset
		id        string
		want      driver.Driver
		wantErr   error
	}{
		{
			name:      "driver configured",
			driverset: driverset,
			id:        "driver0",
			want:      &fake.Driver{FakeID: "driver0"},
			wantErr:   nil,
		},
		{
			name:      "driver not configured",
			driverset: driverset,
			id:        "driver1",
			want:      nil,
			wantErr:   ErrNotConfigured{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.driverset.Get(tc.id)
			assert.IsTypef(t, tc.wantErr, err, "%+#v", tc.driverset)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrDriverDuplicate(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			err := ErrDriverDuplicate{tc.id}
			assert.Equal(t, err.Error(), tc.want)
		})
	}
}

func TestErrNotConfigured(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			err := ErrNotConfigured{tc.id}
			assert.Equal(t, err.Error(), tc.want)
		})
	}
}
