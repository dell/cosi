// Copyright © 2025-2026 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dell/cosi/pkg/config"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

func TestRunMain(t *testing.T) {
	oldRunBlocking := runBlocking
	defer func() { runBlocking = oldRunBlocking }()

	tests := []struct {
		name                   string
		configFile             string
		driverConfigParamsFile string
		otelEndpoint           string
		runBlockingError       error
		resourceOverride       func(ctx context.Context) (*resource.Resource, error)
		clientOverride         func(url string) (*grpc.ClientConn, error)
		traceOverride          func(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error)
		wantErr                bool
	}{
		{
			name:                   "successful start",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "otel-collector:4317",
			runBlockingError:       nil,
			wantErr:                false,
		},
		{
			name:                   "no otel endpoint",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "",
			runBlockingError:       nil,

			wantErr: false,
		},
		{
			name:                   "missing config files",
			configFile:             "test-data-missing.yaml",
			driverConfigParamsFile: "",
			otelEndpoint:           "",
			runBlockingError:       nil,

			wantErr: true,
		},
		{
			name:                   "error while running",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "",
			runBlockingError:       errors.New("error"),
			wantErr:                true,
		},
		{
			name:                   "error creating resource",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "otelEndpoint",
			resourceOverride: func(_ context.Context) (*resource.Resource, error) {
				return nil, errors.New("error")
			},
			wantErr: false,
		},
		{
			name:                   "error creating client",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "otelEndpoint",
			clientOverride: func(_ string) (*grpc.ClientConn, error) {
				return nil, errors.New("error")
			},
			wantErr: false,
		},
		{
			name:                   "error creating trace client",
			configFile:             "test-data.yaml",
			driverConfigParamsFile: "test/driver-config-params.yaml",
			otelEndpoint:           "otelEndpoint",
			traceOverride: func(_ context.Context, _ *grpc.ClientConn) (*otlptrace.Exporter, error) {
				return nil, errors.New("error")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorHandler := &errorHandler{}
			errorHandler.Handle(errors.New("test"))

			configFile = &tt.configFile
			otelEndpoint = &tt.otelEndpoint
			driverConfigParamsFile = &tt.driverConfigParamsFile

			initFlags()
			setOtelLogger()

			if tt.resourceOverride != nil {
				oldResource := newResource
				defer func() { newResource = oldResource }()
				newResource = tt.resourceOverride
			}

			if tt.clientOverride != nil {
				oldClient := newClient
				defer func() { newClient = oldClient }()
				newClient = tt.clientOverride
			}

			if tt.traceOverride != nil {
				oldTrace := newTraceExporter
				defer func() { newTraceExporter = oldTrace }()
				newTraceExporter = tt.traceOverride
			}

			runBlocking = func(_ context.Context, _ *config.ConfigSchemaJson, _ string) error {
				return tt.runBlockingError
			}

			osExit = func() {}

			err := runMain()

			if (err != nil) != tt.wantErr {
				t.Errorf("runMain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// execute run to verify that it does not panic
			run()
		})
	}
}

func TestMainReturnsOnInvalidDriverConfigParams(t *testing.T) {
	configParamsFile := t.TempDir() + "/driver-config-params.yaml"
	if err := os.WriteFile(configParamsFile, []byte("COSI_LOG_LEVEL: INVALID\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	configPath := "test-data.yaml"
	endpoint := ""
	configFile = &configPath
	otelEndpoint = &endpoint
	driverConfigParamsFile = &configParamsFile
	oldRunBlocking := runBlocking
	oldOSExit := osExit
	defer func() {
		runBlocking = oldRunBlocking
		osExit = oldOSExit
	}()
	runBlocking = func(_ context.Context, _ *config.ConfigSchemaJson, _ string) error {
		return nil
	}
	osExit = func() {}

	main()
}

func TestRunMainHandlesDriverConfigChangeError(t *testing.T) {
	configParamsFile := t.TempDir() + "/driver-config-params.yaml"
	if err := os.WriteFile(configParamsFile, []byte("COSI_LOG_LEVEL: INFO\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	configPath := "test-data.yaml"
	endpoint := ""
	configFile = &configPath
	otelEndpoint = &endpoint
	driverConfigParamsFile = &configParamsFile
	oldRunBlocking := runBlocking
	defer func() { runBlocking = oldRunBlocking }()
	runBlocking = func(_ context.Context, _ *config.ConfigSchemaJson, _ string) error {
		if err := os.WriteFile(configParamsFile, []byte("COSI_LOG_LEVEL: INVALID\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	if err := runMain(); err != nil {
		t.Fatal(err)
	}
}

func TestRunBlockingReturnsConfigurationError(t *testing.T) {
	cfg := &config.ConfigSchemaJson{
		Connections: []config.Configuration{{}},
	}

	if err := runBlocking(context.Background(), cfg, ""); err == nil {
		t.Fatal("runBlocking() error = nil, want configuration error")
	}
}

func TestUpdateDriverConfigParams(t *testing.T) {
	tests := []struct {
		name        string
		setLevel    bool
		setFormat   bool
		levelValue  string
		formatValue string
		expectErr   bool
	}{
		{
			name:        "Log level and format unset",
			setLevel:    false,
			setFormat:   false,
			levelValue:  "",
			formatValue: "",
			expectErr:   false,
		},
		{
			name:        "Valid inputs (lowercase)",
			setLevel:    true,
			setFormat:   true,
			levelValue:  "INFO",
			formatValue: "TEXT",
			expectErr:   false,
		},
		{
			name:        "Valid inputs (uppercase)",
			setLevel:    true,
			setFormat:   true,
			levelValue:  "debug",
			formatValue: "json",
			expectErr:   false,
		},
		{
			name:        "Invalid inputs",
			setLevel:    true,
			setFormat:   true,
			levelValue:  "INVALID",
			formatValue: "INVALID",
			expectErr:   true,
		},
		{
			name:        "Use defaults",
			setLevel:    true,
			setFormat:   true,
			levelValue:  "warn",
			formatValue: "",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()

			if tt.setLevel {
				v.Set(ParamLogLevel, tt.levelValue)
			}

			if tt.setFormat {
				v.Set(ParamLogFormat, tt.formatValue)
			}

			err := updateDriverConfigParams(context.Background(), v)
			if (err != nil) != tt.expectErr {
				t.Fatalf("got error %v, expected error %v", err, tt.expectErr)
			}
		})
	}
}
