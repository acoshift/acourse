package internal

import (
	"log"
	"os"

	"cloud.google.com/go/logging"
)

// Loggers
var (
	InfoLogger      = log.New(os.Stderr, "", log.LstdFlags)
	NoticeLogger    = log.New(os.Stderr, "", log.LstdFlags)
	WarningLogger   = log.New(os.Stderr, "", log.LstdFlags)
	ErrorLogger     = log.New(os.Stderr, "", log.LstdFlags)
	EmergencyLogger = log.New(os.Stderr, "", log.LstdFlags)
)

// SetLogger sets loggers
func SetLogger(logger *logging.Logger) {
	InfoLogger = logger.StandardLogger(logging.Info)
	NoticeLogger = logger.StandardLogger(logging.Notice)
	WarningLogger = logger.StandardLogger(logging.Warning)
	ErrorLogger = logger.StandardLogger(logging.Error)
	EmergencyLogger = logger.StandardLogger(logging.Emergency)
}
