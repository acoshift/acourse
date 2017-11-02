package app

import (
	"database/sql"
	"encoding/gob"
)

// Config use to init app package
type Config struct {
	Controller     Controller
	Repository     Repository
	View           View
	DB             *sql.DB
	ProjectID      string
	ServiceAccount []byte
	BucketName     string
	EmailServer    string
	EmailPort      int
	EmailUser      string
	EmailPassword  string
	EmailFrom      string
	BaseURL        string
	XSRFSecret     string
	SQLURL         string
	RedisAddr      string
	RedisPass      string
	RedisPrefix    string
	SessionSecret  []byte
	SlackURL       string
}

func init() {
	gob.Register(sessionKey(0))
}
