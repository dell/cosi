package provisioner

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	driver "github.com/dell/cosi-driver/pkg/provisioner/virtualdriver"
)

// Driverset is a structure holding list of Drivers, that can be added or extracted based on the ID.
type Driverset struct {
	drivers sync.Map
}

// Add is used to add new driver to the Driverset.
func (ds *Driverset) Add(newDriver driver.Driver) error {
	id := newDriver.ID()

	if _, ok := ds.drivers.Load(id); ok {
		return fmt.Errorf("failed to load new driver to driverset sync.Map: %w", ErrDriverDuplicate{id})
	}

	ds.drivers.Store(id, newDriver)

	return nil
}

// Get is used to get driver from the Driverset.
func (ds *Driverset) Get(id string) (driver.Driver, error) {
	d, ok := ds.drivers.Load(id)
	if !ok {
		return nil, fmt.Errorf("failed to get driver from driverset: %w", ErrNotConfigured{id})
	}

	switch d := d.(type) {
	case driver.Driver:
		log.WithFields(log.Fields{
			"id": id,
		}).Debug("driver exists")

		return d, nil
	default:
		return nil, fmt.Errorf("failed to get driver from driverset: %w", errors.New("invalid type"))
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
