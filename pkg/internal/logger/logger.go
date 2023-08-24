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

// Package logger contains interface for logger
// which allows easy switching between logger implementations.
package logger

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
)

const (
	traceLevel = 2
)

// Logger implements aws.Logger interface.
type Logger struct {
	impl logr.Logger
}

var _ aws.Logger = (*Logger)(nil) // interface guard

// New returns new instance of Logger, with logger as implementation.
func New(logger logr.Logger) Logger {
	return Logger{
		impl: logger,
	}
}

// Log is a method that implements aws.Logger interface.
func (l Logger) Log(args ...interface{}) {
	l.impl.V(traceLevel).Info("internal logger message", args...)
}
