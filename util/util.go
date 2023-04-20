package util

import (
	log "github.com/sirupsen/logrus"
)

// SetLogLevel sets the log level based on the logLevel string.
func SetLogLevel(logLevel string) {
	log.SetReportCaller(false)
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
		// SetReportCaller adds the calling method as a field.
		log.SetReportCaller(true)
	case "debug":
		log.SetLevel(log.DebugLevel)
		// SetReportCaller adds the calling method as a field.
		log.SetReportCaller(true)
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
			"log-level":     logLevel,
			"new-log-level": "debug",
		}).Error("unknown log level, setting to debug")
		log.SetLevel(log.DebugLevel)

		return
	}

	log.WithFields(log.Fields{
		"log-level": logLevel,
	}).Info("log level set")
}

func TraceLogging(err error, msg string) error {
	log.WithFields(log.Fields{
		"error_msg": err,
	}).Trace(msg)

	return err
}

func ErrorLogging(err error, msg string) error {
	log.WithFields(log.Fields{
		"error_msg": err,
	}).Error(msg)

	return err
}
