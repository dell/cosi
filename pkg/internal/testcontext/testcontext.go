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

// Package testcontext creates context used in tests for other packages.
package testcontext

import (
	"context"
	"testing"
	"time"
)

const defaultTimeout = 20 * time.Second

// New creates new context with deadline equal to test deadline, or (if the deadline is empty),
// with a timeout equal to the default timeout.
func New(t *testing.T) (context.Context, context.CancelFunc) {
	return getContextCancelFunc(t.Deadline())
}

func getContextCancelFunc(deadline time.Time, ok bool) (context.Context, context.CancelFunc) {
	if ok {
		return context.WithDeadline(context.Background(), deadline)
	}

	return context.WithTimeout(context.Background(), defaultTimeout)
}
