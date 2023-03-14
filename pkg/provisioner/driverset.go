package provisioner

import (
	driver "github.com/dell/cosi-driver/pkg/provisioner/virtual_driver"
)

// Driverset is a structure holding list of Drivers, that can be added or extracted based on the ID
type Driverset struct {
	drivers map[string]driver.Driver
}

// Add is used to add new driver to the Driverset
func (d *Driverset) Add(driver driver.Driver) error {
	id := driver.ID()
	if _, ok := d.drivers[id]; ok == true {
		return ErrDriverDuplicate{id}
	}

	d.drivers[id] = driver

	return nil
}

// Get is used to get driver from the Driverset
func (d Driverset) Get(id string) (driver.Driver, error) {
	driver, ok := d.drivers[id]
	if !ok {
		return nil, ErrNotConfigured{id}
	}

	return driver, nil
}

// ErrDriverDuplicate indicates that the Driver is already present in driverset
type ErrDriverDuplicate struct {
	ID string
}

func (err ErrDriverDuplicate) Error() string {
	return "driver for '" + err.ID + "' already exists"
}

// ErrNotConfigured indicates that the Driver is not present in the driverset
type ErrNotConfigured struct {
	ID string
}

func (err ErrNotConfigured) Error() string {
	return "platform identified by '" + err.ID + "' was not configured"
}
