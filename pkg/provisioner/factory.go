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
	"errors"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtual_driver"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/provisioner/objectscale"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
)

// NewVirtualDriver is factory function, that takes configuration, validates if it is correct, and
// returns correct driver.
func NewVirtualDriver(config config.Configuration) (driver.Driver, error) {
	// in the future, here can be more than one
	if !exactlyOne(config.Objectscale) {
		return nil, errors.New("expected exactly one OSP in configuration")
	}

	if config.Objectscale != nil {
		return objectscale.New(config.Objectscale)
	}

	panic("programming error")
}

func exactlyOne(nillables ...interface{}) bool {
	c := 0
	for _, v := range nillables {
		if v != nil {
			c++
		}
	}

	return c == 1
}
