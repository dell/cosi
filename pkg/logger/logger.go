// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package logger contains interface for logger
// which allows easy switching between logger implementations.
package logger

import (
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

func New() *logrus.Logger {
	logrusInstance := logrus.New()
	logrusInstance.SetReportCaller(false)

	// Set level
	logrusInstance.SetLevel(logrus.Level(defaultLevel))

	// Set formatter
	logrusInstance.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   true, // always print full timestamp
		DisableColors:   true, // never use colors in logs, even if the terminal supports it
	})

	log = logrusr.New(logrusInstance)
	return logrusInstance
}

func Log() logr.Logger {
	if log.GetSink() == nil {
		_ = New()
	}
	return log
}
