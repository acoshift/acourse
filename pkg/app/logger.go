package app

import (
	"log"

	"cloud.google.com/go/logging"
	"github.com/acoshift/acourse/pkg/internal"
)

// SetLogger sets internal loggers
func SetLogger(logger *logging.Logger) {
	internal.SetLogger(logger)
}

// NewNoFatalLogger creates new logger
// TODO: This logger is a heck for grpc fatal
// this shouls be removed after grpc can handle `sendResponse` with nil response
// without calling grpclog.Fatal
func NewNoFatalLogger() *NoFatalLogger {
	return &NoFatalLogger{}
}

// NoFatalLogger implements groclog.Logger
type NoFatalLogger struct{}

// Fatal prevents server fatal by using Print instead
func (*NoFatalLogger) Fatal(args ...interface{}) {
	log.Print(args...)
}

// Fatalf prevents server fatal by using Printf instead
func (*NoFatalLogger) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Fatalln prevents server fatal by using Println instead
func (*NoFatalLogger) Fatalln(args ...interface{}) {
	log.Println(args...)
}

// Print wraps log.Print
func (*NoFatalLogger) Print(args ...interface{}) {
	log.Print(args...)
}

// Printf wraps log.Printf
func (*NoFatalLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Println wraps log.Println
func (*NoFatalLogger) Println(args ...interface{}) {
	log.Println(args...)
}
