package app

import (
	"log"
)

// NewLogger creates new logger
// TODO: This logger is a heck for grpc fatal
// this shouls be removed after grpc can handle `sendResponse` with nil response
// without calling grpclog.Fatal
func NewLogger() *Logger {
	return &Logger{}
}

// Logger implements groclog.Logger
type Logger struct{}

// Fatal prevents server fatal by using Print instead
func (*Logger) Fatal(args ...interface{}) {
	log.Print(args...)
}

// Fatalf prevents server fatal by using Printf instead
func (*Logger) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Fatalln prevents server fatal by using Println instead
func (*Logger) Fatalln(args ...interface{}) {
	log.Println(args...)
}

// Print wraps log.Print
func (*Logger) Print(args ...interface{}) {
	log.Print(args...)
}

// Printf wraps log.Printf
func (*Logger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Println wraps log.Println
func (*Logger) Println(args ...interface{}) {
	log.Println(args...)
}
