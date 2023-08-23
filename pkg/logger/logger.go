package logger

import (
	"github.com/bombsimon/logrusr/v4"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

const (
	WarnLevel = iota - 1 // This will get ignored because it's below InfoLevel V(0), we don't have access to logrus level 0 - 3 because of logrusr implementation, 3 being logrus.WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger struct {
	Logger logr.Logger
}

var Log *Logger

func New(level string, formatter string) {
	logrusInstance := logrus.New()
	logrusInstance.SetReportCaller(false)

	// Set level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrusInstance.SetLevel(logrus.DebugLevel)
	} else {
		logrusInstance.SetLevel(lvl)
	}

	// Set formatter
	timestampFormat := "2006-01-02 15:04:05.000"

	switch formatter {
	case "json":
		logrusInstance.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     false, // do not indent JSON logs, print each log entry on one line
		})

	case "text":
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})

	case "pretty":
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   false, // do not print full timestamps
			DisableColors:   false, // do not disable colors
		})

	default:
		logrusInstance.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: timestampFormat,
			FullTimestamp:   true, // always print full timestamp
			DisableColors:   true, // never use colors in logs, even if the terminal supports it
		})
	}

	logrInstance := logrusr.New(logrusInstance)

	Log = &Logger{
		Logger: logrInstance,
	}
}

func GetLogger() *Logger {
	return Log
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.Logger.V(DebugLevel).Info(msg, keysAndValues)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info(msg, keysAndValues)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	// l.Logger.V(WarnLevel).Info(fmt.Sprintf("WARN: %s", msg), keysAndValues)
	l.Logger.V(WarnLevel).Info(msg, keysAndValues)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.Logger.Error(err, msg, keysAndValues)
}

func (l *Logger) Trace(msg string, keysAndValues ...interface{}) {
	l.Logger.V(TraceLevel).Info(msg, keysAndValues)
}
