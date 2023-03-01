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
	"net"
	"testing"
	"time"
)

// FIXME: some way to test this? probaby refactor is needed
func TestRun_Successful(t *testing.T) {
	// Test server starts successfully and stops gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(ctx, "test", 8080)
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Cancel context to stop server gracefully
	cancel()

	err := <-errCh
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRun_PortAlreadyInUse(t *testing.T) {
	// Test error is returned when port is already in use
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("Failed to start test listener: %v", err)
	}
	defer lis.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(ctx, "test", 8080)
	}()

	err = <-errCh
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}
