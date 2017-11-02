package app

import (
	"database/sql"
	"encoding/gob"
)

// Config use to init app package
type Config struct {
	Controller    Controller
	Repository    Repository
	View          View
	DB            *sql.DB
	BaseURL       string
	XSRFSecret    string
	RedisAddr     string
	RedisPass     string
	RedisPrefix   string
	SessionSecret []byte
}

func init() {
	gob.Register(sessionKey(0))
}
