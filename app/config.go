package app

import (
	"database/sql"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/go-redis/redis"
	"gopkg.in/gomail.v2"
)

// Config use to init app package
type Config struct {
	DB           *sql.DB
	BaseURL      string
	RedisClient  *redis.Client
	RedisPrefix  string
	Auth         *firebase.Auth
	Location     *time.Location
	SlackURL     string
	EmailFrom    string
	EmailDialer  *gomail.Dialer
	BucketHandle *storage.BucketHandle
	BucketName   string
}
