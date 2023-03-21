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
	driver "github.com/dell/cosi-driver/pkg/provisioner/virtual_driver"
	"github.com/dell/cosi-driver/pkg/provisioner/virtual_driver/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	driverset = Driverset{
		drivers: map[string]driver.Driver{
			"driver0": &fake.Driver{FakeId: "driver0"},
		}}
)

func TestDriversetAdd(t *testing.T) {
	testCases := []struct {
		name      string
		driverset Driverset
		driver    fake.Driver
		want      Driverset
		wantErr   error
	}{
		{
			name:      "no duplicate",
			driverset: Driverset{drivers: map[string]driver.Driver{}},
			driver:    fake.Driver{FakeId: "driver0"},
			want:      driverset,
			wantErr:   nil,
		},
		{
			name:      "duplicate",
			driverset: driverset,
			driver:    fake.Driver{FakeId: "driver0"},
			want:      driverset,
			wantErr:   ErrDriverDuplicate{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.driverset.Add(&tc.driver)
			assert.IsType(t, tc.wantErr, err)
			assert.Equal(t, tc.want.drivers, tc.driverset.drivers)
		})
	}
}

func TestDriversetGet(t *testing.T) {
	testCases := []struct {
		name      string
		driverset Driverset
		id        string
		want      driver.Driver
		wantErr   error
	}{
		{
			name:      "driver configured",
			driverset: driverset,
			id:        "driver0",
			want:      &fake.Driver{FakeId: "driver0"},
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
