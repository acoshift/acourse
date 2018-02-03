package app

import (
	"database/sql"
	"encoding/gob"

	"github.com/garyburd/redigo/redis"
)

// Config use to init app package
type Config struct {
	Controller    Controller
	Repository    Repository
	DB            *sql.DB
	BaseURL       string
	XSRFSecret    string
	RedisPool     *redis.Pool
	RedisPrefix   string
	CachePool     *redis.Pool
	CachePrefix   string
	SessionSecret []byte
}

func init() {
	gob.Register(sessionKey(0))
}
