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
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/pkg/driver"
	"github.com/dell/cosi-driver/util"
)

var (
	logLevel   = flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	configFile = flag.String("config", "/cosi/config.yaml", "path to config file")
)

const (
	tracedServiceName = "cosi-driver"
	jaegerURL         = "http://10.247.103.53/collector/"
)

// init is run before main and is used to define command line flags.
func init() {
	// Parse command line flags.
	flag.Parse()
	// Set the log level.
	util.SetLogLevel(*logLevel)
	// Set the log format.
	util.SetLoggingFormatter()
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
		"config_file_path": configFile,
	}).Info("config successfully loaded")

	// Create TracerProvider with exporter to Jaeger.
	// TODO: let user configure jaeger url.
	tp, err := tracerProvider(jaegerURL)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("failed to connect to Jaeger")
	}
	// Set global TracerProvider.
	otel.SetTracerProvider(tp)

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
	return driver.RunBlocking(ctx, cfg, driver.COSISocket, "cosi-driver")
}

// tracerProvider creates new tracerProvider and connects it to Jaeger running under provided URL.
func tracerProvider(url string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(tracedServiceName),
		)),
	)

	return tp, nil
}
