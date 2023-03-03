package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi-driver/core"
)

// SetLogLevel sets the log level based on the logLevel string
func SetLogLevel(logLevel string) {
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
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
	}
}

// PrintVersion
func PrintVersion() {
	fmt.Println(core.SemVer)
}
