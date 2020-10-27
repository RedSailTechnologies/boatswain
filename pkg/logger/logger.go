package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var l *zap.SugaredLogger
var once sync.Once

func init() {
	initializeLogger()
}

// Debug log level
func Debug(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Debugw(message, args...)
	} else {
		l.Debug(message)
	}
}

// Info log level
func Info(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Infow(message, args...)
	} else {
		l.Info(message)
	}
}

// Warn log level
func Warn(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Warnw(message, args...)
	} else {
		l.Warn(message)
	}
}

// Error log level
func Error(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Errorw(message, args...)
	} else {
		l.Error(message)
	}
}

// Panic log level
func Panic(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Panicw(message, args...)
	} else {
		l.Panic(message)
	}
}

// Fatal log level
func Fatal(message string, args ...interface{}) {
	if len(args) > 0 {
		l.Fatalw(message, args...)
	} else {
		l.Fatal(message)
	}
}

func initializeLogger() {
	once.Do(func() {
		z, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("error creating logger")
		}
		l = z.Sugar()
	})
}
