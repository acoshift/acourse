package app

import (
	"database/sql"
	"encoding/gob"
	"time"

	"cloud.google.com/go/storage"
	firebase "github.com/acoshift/go-firebase-admin"
	"github.com/garyburd/redigo/redis"
	gomail "gopkg.in/gomail.v2"
)

// Config use to init app package
type Config struct {
	DB            *sql.DB
	BaseURL       string
	RedisPool     *redis.Pool
	RedisPrefix   string
	CachePool     *redis.Pool
	CachePrefix   string
	SessionSecret []byte
	Auth          *firebase.Auth
	Location      *time.Location
	SlackURL      string
	EmailFrom     string
	EmailDialer   *gomail.Dialer
	BucketHandle  *storage.BucketHandle
	BucketName    string
}

func init() {
	gob.Register(sessionKey(0))
}
