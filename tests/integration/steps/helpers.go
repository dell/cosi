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

package steps

import (
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
)

const (
	attempts = 5
	sleep    = 2 * time.Second // nolint:gomnd
)

func retry(ctx ginkgo.SpecContext, attempts int, sleep time.Duration, f func() error) error {
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
