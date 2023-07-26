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

	log "github.com/sirupsen/logrus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/bombsimon/logrusr/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/driver"
)

var (
	logLevel     = flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
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
	// Set the log format.
	// This must be done before the log level is set, so if any errors occur, they are logged in proper format.
	setLogFormatter(*logFormat)
	// Set the log level.
	setLogLevel(*logLevel)
	// Set the custom logger for OpenTelemetry.
	setOtelLogger()
}

func main() {
	err := runMain()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to start application")
	}
}

func runMain() error {
	// Create a context that is canceled when the SIGINT or SIGTERM signal is received.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to create configuration")
	}

	log.WithFields(log.Fields{
		"configFilePath": *configFile,
	}).Info("config successfully loaded")

	// Create TracerProvider with exporter to Open Telemetry Collector.
	var tp *sdktrace.TracerProvider
	if *otelEndpoint != "" {
		tp, err = tracerProvider(ctx, *otelEndpoint)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("failed to connect to Jaeger")
		} else {
			// Set global TracerProvider.
			otel.SetTracerProvider(tp)
			// set global propagator to tracecontext (the default is no-op).
			otel.SetTextMapPropagator(propagation.TraceContext{})
			log.WithFields(log.Fields{
				"collector": *otelEndpoint,
			}).Info("tracing started successfully")
		}
	} else {
		log.Warn("OTEL endpoint is empty, disabling tracing; please refer to helm's values.yaml")
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
		log.WithFields(log.Fields{
			"type": sig,
		}).Info("signal received")
		// Cancel the context.
		cancel()
		// Exit the program with an error.
		os.Exit(1)
	}()

	log.Info("COSI driver is starting")
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

// setLogLevel sets the log level based on the logLevel string.
func setLogLevel(logLevel string) {
	log.SetReportCaller(false)

	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
		// SetReportCaller adds the calling method as a field.
		log.SetReportCaller(true)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.WithFields(log.Fields{
			"logLevel":    logLevel,
			"newLogLevel": "debug",
		}).Error("unknown log level, setting to debug")
		log.SetLevel(log.DebugLevel)

		return
	}

	log.WithFields(log.Fields{
		"logLevel": logLevel,
	}).Info("log level set")
}

// setLogFormatter set is used to set proper formatter for logs.
func setLogFormatter(logFormat string) {
	timestampFormat := "2006-01-02 15:04:05.000"

	switch logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     false, // do not indent JSON logs, print each log entry on one line
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})

	case "pretty":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   false, // do not print full timestamps
			DisableColors:   false, // do not disable colors
		})

	default:
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})

		log.WithFields(log.Fields{
			"logFormat":    logFormat,
			"newLogFormat": "text",
		}).Error("unknown log format, setting to text")
	}
}

// errorHandler implements otel.ErrorHandler interface.
type errorHandler struct{}

// Handle is used to handle errors from OpenTelemetry.
func (e *errorHandler) Handle(err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("error occurred in OpenTelemetry")
}

// setOtelLogger is used to set the custom logger from OpenTelemetry.
func setOtelLogger() {
	logger := logrusr.New(log.StandardLogger())
	otel.SetLogger(logger)
	otel.SetErrorHandler(&errorHandler{})
}
