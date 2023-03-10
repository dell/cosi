package virtual_driver

import (
	"errors"

	"github.com/dell/cosi-driver/pkg/config"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
)

// New is factory function, that takes configuration, validates if it is correct, and
// returns correct driver.
func New(config config.Configuration) (Driver, error) {
	if !exactlyOne(config.Ecs, config.Objectscale, config.Powerscale) {
		return nil, errors.New("expected exactly one OSP in configuration")
	}

	if config.Ecs != nil {
		return newEcsDriver(config.Ecs)
	} else if config.Objectscale != nil {
		return newObjectscaleDriver(config.Objectscale)
	} else if config.Powerscale != nil {
		return newPowerscaleDriver(config.Powerscale)
	}

	panic("programming error")
}

// newEcsDriver is function that takes Ecs configuration and builds
// correct driver based on it's specification. This should configure
// TLS, managment client, S3 endpoint and all other necessary fields
func newEcsDriver(config *config.Ecs) (Driver, error) {
	return nil, ErrUnimplemented
}

// newObjectscaleDriver is function that takes Objectscale configuration
// and builds correct driver based on it's specification. This should
// configure TLS, managment client, IAM client, S3 endpoint and all
// other necessary fields
func newObjectscaleDriver(config *config.Objectscale) (Driver, error) {
	return nil, ErrUnimplemented
}

// newPowerstoreDriver is function that takes Powerstore configuration
// and builds correct driver based on it's specification. This should
// configure TLS, managment client, S3 endpoint and all other
// necessary fields
func newPowerscaleDriver(config *config.Powerscale) (Driver, error) {
	return nil, ErrUnimplemented
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
