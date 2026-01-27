// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package provisioner

import (
	"errors"
	"fmt"
	"sync"

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
