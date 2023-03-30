package provisioner

import (
	"sync"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
)

// Driverset is a structure holding list of Drivers, that can be added or extracted based on the ID
type Driverset struct {
	once    sync.Once
	drivers map[string]driver.Driver
}

func (ds *Driverset) init() {
	if ds.drivers == nil {
		ds.drivers = make(map[string]driver.Driver)
	}
}

// Add is used to add new driver to the Driverset
func (ds *Driverset) Add(newDriver driver.Driver) error {
	id := newDriver.ID()
	if _, ok := ds.drivers[id]; ok {
		return ErrDriverDuplicate{id}
	}

	ds.once.Do(ds.init)
	ds.drivers[id] = newDriver

	return nil
}

// Get is used to get driver from the Driverset
func (ds *Driverset) Get(id string) (driver.Driver, error) {
	ds.once.Do(ds.init)

	driver, ok := ds.drivers[id]
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
