package internal

import (
	"log"

	"cloud.google.com/go/logging"
)

// Loggers
var (
	InfoLogger      *log.Logger
	NoticeLogger    *log.Logger
	WarningLogger   *log.Logger
	ErrorLogger     *log.Logger
	EmergencyLogger *log.Logger
)

// SetLogger sets loggers
func SetLogger(logger *logging.Logger) {
	InfoLogger = logger.StandardLogger(logging.Info)
	NoticeLogger = logger.StandardLogger(logging.Notice)
	WarningLogger = logger.StandardLogger(logging.Warning)
	ErrorLogger = logger.StandardLogger(logging.Error)
	EmergencyLogger = logger.StandardLogger(logging.Emergency)
}
