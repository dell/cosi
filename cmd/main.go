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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/driver"
	logger "github.com/dell/cosi/pkg/logger"
	"github.com/go-logr/logr"
)

var log logr.Logger

var (
	logLevel     = flag.Int("log-level", 4, "Log level (0-10)")
	logFormat    = flag.String("log-format", "text", "Log format (text, json)")
	otelEndpoint = flag.String("otel-endpoint", "",
		"OTEL collector endpoint for collecting observability data")
	configFile = flag.String("config", "/cosi/config.yaml", "path to config file")
)

const (
	tracedServiceName = "cosi.dellemc.com"
)

// init is run before main and is used to define command line flags.
func init() {
	// Parse command line flags.
	flag.Parse()
	// Create logger instance.
	logger.New(*logLevel, *logFormat)
	log = logger.GetLogger()
	// Set the custom logger for OpenTelemetry.
	setOtelLogger()
}

func main() {
	err := runMain()
	if err != nil {
		log.Error(err, "failed to start application")
		os.Exit(1)
	}
}

func runMain() error {
	// Create a context that is canceled when the SIGINT or SIGTERM signal is received.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(*configFile)
	if err != nil {
		log.Error(err, "failed to create configuration")
		return err
	}

	log.V(4).Info("Config successfully loaded", "configFilePath", *configFile)

	// Create TracerProvider with exporter to Open Telemetry Collector.
	var tp *sdktrace.TracerProvider
	if *otelEndpoint != "" {
		tp, err = tracerProvider(ctx, *otelEndpoint)
		if err != nil {
			log.V(0).Info("Failed to connect to Jaeger", "error", err)
		} else {
			// Set global TracerProvider.
			otel.SetTracerProvider(tp)
			// set global propagator to tracecontext (the default is no-op).
			otel.SetTextMapPropagator(propagation.TraceContext{})
			log.V(4).Info("Tracing started successfully", "collector", *otelEndpoint)
		}
	} else {
		log.V(0).Info("OTEL endpoint is empty, disabling tracing")
	}

	// Create a channel to listen for signals.
	sigs := make(chan os.Signal, 1)
	// Listen for the SIGINT and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a goroutine to listen for signals.
	go func() {
		// Wait for a signal.
		sig := <-sigs
		// Log that a signal was received.
		log.V(4).Info("Signal received", "type", sig)
		// Cancel the context.
		cancel()
		// Exit the program with an error.
		os.Exit(1)
	}()

	log.V(4).Info("COSI driver is starting")
	// Run the driver.
	return driver.RunBlocking(ctx, cfg, driver.COSISocket, tracedServiceName)
}

// tracerProvider creates new tracerProvider and connects it to Jaeger running under provided URL.
func tracerProvider(ctx context.Context, url string) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(tracedServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracing resource: %w", err)
	}

	conn, err := grpc.DialContext(
		ctx,
		url,
		// insecure transport left intentionally here
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tp, nil
}

// errorHandler implements otel.ErrorHandler interface.
type errorHandler struct{}

// Handle is used to handle errors from OpenTelemetry.
func (e *errorHandler) Handle(err error) {
	log.Error(err, "error occurred in OpenTelemetry")
}

// setOtelLogger is used to set the custom logger from OpenTelemetry.
func setOtelLogger() {
	otel.SetLogger(log)
	otel.SetErrorHandler(&errorHandler{})
}
