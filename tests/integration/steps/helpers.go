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
	"fmt"
	"time"
)

const (
	attempts = 60
	sleep    = 2 * time.Second //nolint:gomnd
)

func retry(ctx context.Context, attempts int, sleep time.Duration, f func() error) error {
	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	retries := 0

	for {
		select {
		case <-ticker.C:
			retries++

			err := f()
			if err == nil {
				if retries > 1 {
					fmt.Printf("[retry] succeeded on attempt %d/%d\n", retries, attempts)
				}

				return nil
			}

			fmt.Printf("[retry] attempt %d/%d failed: %v\n", retries, attempts, err)

			if retries >= attempts {
				fmt.Printf("[retry] exhausted all %d attempts\n", attempts)
				return err
			}

		case <-ctx.Done():
			fmt.Printf("[retry] context cancelled after %d attempts: %v\n", retries, ctx.Err())
			return ctx.Err()
		}
	}
}
