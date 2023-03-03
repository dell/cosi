//Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
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

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi-driver/pkg/driver"
	"github.com/dell/cosi-driver/util"
)

var (
	versionFlag = flag.Bool("version", false, "Print the version and exit.")
	logLevel    = flag.String("log-level", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	port        = flag.Int("port", 9000, "Port to listen on")
)

// init is run before main and is used to define command line flags.
func init() {
	// Set standard logger output.
	log.SetOutput(os.Stdout)
	// Parse command line flags.
	flag.Parse()
	// Set the log level.
	util.SetLogLevel(*logLevel)
	// If the version flag is specified, print the version and exit.
	if *versionFlag {
		util.PrintVersion()
		os.Exit(0)
	}
}

func main() {
	// Create a context that is canceled when the SIGINT or SIGTERM signal is received.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for signals.
	sigs := make(chan os.Signal, 1)
	// Listen for the SIGINT and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a goroutine to listen for signals.
	go func() {
		// Wait for a signal.
		sig := <-sigs
		// Log that a signal was received.
		log.Info("Signal received", "type", sig)
		// Cancel the context.
		cancel()
		// Exit the program with an error.
		os.Exit(1)
	}()

	// Run the driver.
	err := driver.Run(ctx, "cosi-driver", *port)
	if err != nil {
		log.Fatal(err)
	}
}
