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

type Logger struct{}

func (*Logger) Fatal(args ...interface{}) {
	log.Print(args...)
}

func (*Logger) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (*Logger) Fatalln(args ...interface{}) {
	log.Println(args...)
}

func (*Logger) Print(args ...interface{}) {
	log.Print(args...)
}

func (*Logger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (*Logger) Println(args ...interface{}) {
	log.Println(args...)
}
