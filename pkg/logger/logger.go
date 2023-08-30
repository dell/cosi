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
	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

const (
	minLevel     = 0
	maxLevel     = 10
	defaultLevel = 4

	timestampFormat = "2006-01-02 15:04:05.000"
)

// AWSLogger implements aws.Logger interface.
type AWSLogger struct {
	impl logr.Logger
}

var _ aws.Logger = (*AWSLogger)(nil) // interface guard

// NewAWSLogger returns new instance of Logger, with logger as implementation.
func NewAWSLogger(logger logr.Logger) AWSLogger {
	return AWSLogger{
		impl: logger,
	}
}

// Log is a method that implements aws.Logger interface.
func (l AWSLogger) Log(keysAndValues ...interface{}) {
	l.impl.V(2).Info("internal logger message", keysAndValues...)
}

var log logr.Logger

func New(level int, formatter string) {
	logrusInstance := logrus.New()
	logrusInstance.SetReportCaller(false)

	// Set level
	logrusInstance.SetLevel(logrus.Level(defaultLevel))

	if level >= minLevel || level <= maxLevel {
		logrusInstance.SetLevel(logrus.Level(level))
	}

	switch formatter {
	case "json":
		logrusInstance.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     false, // do not indent JSON logs, print each log entry on one line
		})

	case "text":
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})

	case "pretty":
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   false, // do not print full timestamps
			DisableColors:   false, // do not disable colors
		})

	default:
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})
	}

	log = logrusr.New(logrusInstance)
}

func GetLogger() logr.Logger {
	return log
}
