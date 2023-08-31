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
	"fmt"
	"sync"

	l "github.com/dell/cosi/pkg/logger"
	driver "github.com/dell/cosi/pkg/provisioner/virtualdriver"
)

// Driverset is a structure holding list of Drivers, that can be added or extracted based on the ID.
type Driverset struct {
	drivers sync.Map
}

// Add is used to add new driver to the Driverset.
func (ds *Driverset) Add(newDriver driver.Driver) error {
	id := newDriver.ID()

	if _, ok := ds.drivers.Load(id); ok {
		return fmt.Errorf("failed to load new configuration for specified object storage platform: %w", ErrDriverDuplicate{id})
	}

	ds.drivers.Store(id, newDriver)

	return nil
}

// Get is used to get driver from the Driverset.
func (ds *Driverset) Get(id string) (driver.Driver, error) {
	d, ok := ds.drivers.Load(id)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve configuration for specified object storage platform: %w", ErrNotConfigured{id})
	}

	switch d := d.(type) {
	case driver.Driver:
		l.Log().V(6).Info("Driver exists.", "id", id)

		return d, nil
	default:
		return nil, fmt.Errorf("failed to retrieve configuration for specified object storage platform: %w", errors.New("invalid type"))
	}
}

// ErrDriverDuplicate indicates that the Driver is already present in driverset.
type ErrDriverDuplicate struct {
	ID string
}

func (err ErrDriverDuplicate) Error() string {
	return "driver for '" + err.ID + "' already exists"
}

// ErrNotConfigured indicates that the Driver is not present in the driverset.
type ErrNotConfigured struct {
	ID string
}

func (err ErrNotConfigured) Error() string {
	return "platform identified by '" + err.ID + "' was not configured"
}
