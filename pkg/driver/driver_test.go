//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

package driver

import (
	"context"
	"testing"
	"time"

	"github.com/dell/cosi-driver/pkg/config"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name          string
		port          int
		backendID     string
		namespace     string
		expectedError bool
	}{
		{
			name:          "Successful",
			port:          8090,
			backendID:     "123",
			namespace:     "namespace1",
			expectedError: false,
		},
		{
			name:          "PortAlreadyInUse",
			port:          8090,
			backendID:     "123",
			namespace:     "namespace1",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			// Test server starts successfully and stops gracefully
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			errCh := make(chan error, 1)
			go func() {
				errCh <- Run(ctx, &config.ConfigSchemaJson{}, "test") // FIXME: config is not provided, this will fail!
			}()

			// Wait for server to start
			time.Sleep(500 * time.Millisecond)

			if tc.expectedError {
				// Test error is returned when port is already in use
				err = Run(context.Background(), &config.ConfigSchemaJson{}, "test") // FIXME: config is not provided, this will fail!
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				// Cancel context to stop server gracefully
				cancel()

				err = <-errCh
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
