// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	logger "github.com/dell/cosi/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dell/cosi/pkg/config"
	"github.com/dell/cosi/pkg/driver"
	"github.com/dell/csmlog"
)

var (
	otelEndpoint           = flag.String("otel-endpoint", "", "OTEL collector endpoint for collecting observability data")
	configFile             = flag.String("config", "/cosi/config.yaml", "path to config file")
	driverConfigParamsFile = flag.String("driver-config-params", "", "path to driver config params file")
)

const (
	tracedServiceName = "cosi.dellemc.com"
	// DefaultLogLevel for logs
	DefaultLogLevel = csmlog.InfoLevel
	// ParamLogLevel driver log level
	ParamLogLevel = "COSI_LOG_LEVEL"
	// ParamLogFormat driver log format
	ParamLogFormat = "COSI_LOG_FORMAT"
)

var osExit = func() {
	os.Exit(1)
}

var log = csmlog.GetLogger()

func initFlags() {
	// Parse command line flags.
	flag.Parse()
	// Set the custom logger for OpenTelemetry.
	setOtelLogger()
}

func main() {
	run()
}

func run() {
	initFlags()
	err := runMain()
	if err != nil {
		log.Errorf("failed to start application: %v", err)
		osExit()
	}
}

func runMain() error {
	// Create a context that is canceled when the SIGINT or SIGTERM signal is received.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(*configFile)
	if err != nil {
		return fmt.Errorf("failed to create configuration: %w", err)
	}

	log.Infof("Config successfully loaded from %s", *configFile)

	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(*driverConfigParamsFile)
	if err := v.ReadInConfig(); err != nil {
		log.Warnf("unable to read driver config params from file, using defaults")
	}

	if err := updateDriverConfigParams(ctx, v); err != nil {
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(fsnotify.Event) {
		log.Infof("Driver config params file changed")
		if err := updateDriverConfigParams(ctx, v); err != nil {
			log.Warn(err.Error())
		}
	})

	// Create TracerProvider with exporter to Open Telemetry Collector.
	var tp *sdktrace.TracerProvider
	if *otelEndpoint != "" {
		tp, err = tracerProvider(ctx, *otelEndpoint)
		if err != nil {
			log.Errorf("failed to connect to Jaeger: %v", err)
		} else {
			// Set global TracerProvider.
			otel.SetTracerProvider(tp)
			// set global propagator to tracecontext (the default is no-op).
			otel.SetTextMapPropagator(propagation.TraceContext{})
			log.Infof("Tracing started successfully to %s", *otelEndpoint)
		}
	} else {
		log.Info("OTEL endpoint is empty, disabling tracing")
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
		log.Infof("Signal received: %v", sig)
		// Cancel the context.
		cancel()
		// Exit the program with an error.
		os.Exit(1)
	}()

	log.Info("COSI driver is starting")
	// Run the driver.
	return runBlocking(ctx, cfg, tracedServiceName)
}

func updateDriverConfigParams(ctx context.Context, v *viper.Viper) error {
	log := log.WithContext(ctx)
	logLevel := DefaultLogLevel
	if v.IsSet(ParamLogLevel) {
		inputLogLevel := v.GetString(ParamLogLevel)
		if inputLogLevel != "" {
			parsedLogLevel, err := csmlog.ParseLevel(inputLogLevel)
			if err != nil {
				return fmt.Errorf("invalid log level %s: %v", inputLogLevel, err)
			}

			logLevel = parsedLogLevel
		}
	}

	csmlog.SetLevel(logLevel)
	log.Infof("Log level set to %s", logLevel)

	if v.IsSet(ParamLogFormat) {
		logFormat := strings.ToUpper(v.GetString(ParamLogFormat))
		switch logFormat {
		case "TEXT":
			log.SetFormatter(&logrus.TextFormatter{})
		case "JSON":
			log.SetFormatter(&logrus.JSONFormatter{})
		default:
			// use text formatter by default
			log.SetFormatter(&logrus.TextFormatter{})
		}
		log.Infof("Log format set to %s", logFormat)
	}

	return nil
}

var runBlocking = func(ctx context.Context, cfg *config.ConfigSchemaJson, tracedServiceName string) error {
	return driver.RunBlocking(ctx, cfg, driver.COSISocket, tracedServiceName)
}

var newResource = func(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(tracedServiceName),
		),
	)
}

var newClient = func(url string) (*grpc.ClientConn, error) {
	return grpc.NewClient(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

var newTraceExporter = func(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

// tracerProvider creates new tracerProvider and connects it to Jaeger running under provided URL.
func tracerProvider(ctx context.Context, url string) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := newResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracing resource: %w", err)
	}

	conn, err := newClient(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := newTraceExporter(ctx, conn)
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
	log.Errorf("error occurred in OpenTelemetry: %v", err)
}

// setOtelLogger is used to set the custom logger from OpenTelemetry.
func setOtelLogger() {
	otel.SetLogger(logger.Log())
	otel.SetErrorHandler(&errorHandler{})
}
