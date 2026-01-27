// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package steps

import (
	"context"
	"time"
)

const (
	attempts = 10
	sleep    = 2 * time.Second //nolint:gomnd
)

func retry(ctx context.Context, attempts int, sleep time.Duration, f func() error) error {
	ticker := time.NewTicker(sleep)
	retries := 0

	for {
		select {
		case <-ticker.C:
			err := f()
			if err == nil {
				return nil
			}

			retries++
			if retries > attempts {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
