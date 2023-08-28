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
	"errors"

	driver "github.com/dell/cosi/pkg/provisioner/virtualdriver"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/provisioner/objectscale"
)

// NewVirtualDriver is factory function, that takes configuration, validates if it is correct, and
// returns correct driver.
func NewVirtualDriver(config config.Configuration) (driver.Driver, error) {
	// in the future, here can be more than one
	if !exactlyOne(config.Objectscale) {
		return nil, errors.New("expected exactly one object storage platform in configuration")
	}

	if config.Objectscale != nil {
		log.V(6).Info("ObjectScale config created.")
		return objectscale.New(config.Objectscale)
	}

	panic("programming error")
}

// exactlyOne checks if exactly one of its arguments is not nil.
//
// It takes in a variadic argument list of nillable values (interfaces that can either be nil or non-nil).
// The function then counts the number of non-nil values and returns a boolean indicating whether
// exactly one non-nil value was passed as an argument.
func exactlyOne(nillables ...interface{}) bool {
	count := 0

	for _, nillable := range nillables {
		// we need type switch, because nil not always equals nil, e.g.: `nil != (*config.Objectscale)(nil)`
		switch nillable := nillable.(type) {
		case *config.Objectscale:
			if nillable != (*config.Objectscale)(nil) {
				count++
			}

		default:
			if nillable != nil {
				count++
			}
		}
	}

	return count == 1
}
