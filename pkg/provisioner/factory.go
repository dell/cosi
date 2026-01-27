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

	driver "github.com/dell/cosi/pkg/provisioner/virtualdriver"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/provisioner/objectscale"
)

// NewVirtualDriver is factory function, that takes configuration, validates if it is correct, and
// returns correct driver.
func NewVirtualDriver(config config.Configuration) (driver.Driver, error) {
	// in the future, here can be more than one
	// if !exactlyOne(config.Objectscale) {
	// 	return nil, errors.New("expected exactly one object storage platform in configuration")
	// }

	if config.Objectscale != nil {
		log.Info("ObjectScale config created")
		return objectscale.New(config.Objectscale)
	}

	return nil, errors.New("configuration is empty")
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
