package util

import (
	"context"
	"time"

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
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)
}

func Retry(ctx context.Context, attempts int, sleep time.Duration, f func() error) error {
	ticker := time.NewTicker(sleep)
	retries := 0

	for {
		select {
		case <-ticker.C:
			err := f()
			if err == nil {
				return nil
			}

			retries++
			if retries > attempts {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
