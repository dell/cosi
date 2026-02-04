// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
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
	"testing"

	"github.com/bombsimon/logrusr/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSLogger(t *testing.T) {
	tmpLogger := logrus.New()
	logger := NewAWSLogger(logrusr.New(tmpLogger))
	assert.NotNil(t, logger)
}

func TestAWSLoggerLog(_ *testing.T) {
	tmpLogger := logrus.New()
	logger := NewAWSLogger(logrusr.New(tmpLogger))
	logger.Log("key", "value")
}

func TestNew(t *testing.T) {
	l := New()

	wantLevel := logrus.Level(defaultLevel)
	if logrus.GetLevel() != wantLevel {
		t.Errorf("expected log level %v, got %v", wantLevel, logrus.GetLevel())
	}

	wantFormatter := &logrus.TextFormatter{
		TimestampFormat: timestampFormat,
		FullTimestamp:   true, // always print full timestamp
		DisableColors:   true, // never use colors in logs, even if the terminal supports it
	}
	_, ok := l.Formatter.(*logrus.TextFormatter)
	if !ok {
		t.Errorf("expected log formatter %v, got %v", wantFormatter, l.Formatter)
	}

	_log := Log()
	assert.NotNil(t, _log)
}
