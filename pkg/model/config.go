package model

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
)

// model shared vars
var (
	db          *sql.DB
	redisPool   *redis.Pool
	redisPrefix string
)

// Config use to init model package
type Config struct {
	DB          *sql.DB
	RedisPool   *redis.Pool
	RedisPrefix string
}

// Init inits model package
func Init(config Config) error {
	db = config.DB
	redisPool = config.RedisPool
	redisPrefix = config.RedisPrefix

	return nil
}
