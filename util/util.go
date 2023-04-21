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

// SetLoggingFormatter set timestamp in logs.
func SetLoggingFormatter() {
	formatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)
}

// ErrorLogging log error and message where it failed.
func ErrorLogging(err error, msg string) error {
	log.WithFields(log.Fields{
		"error": err,
	}).Error(msg)

	return err
}
