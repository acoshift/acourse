package model

import "github.com/garyburd/redigo/redis"

// model shared vars
var (
	redisPool   *redis.Pool
	redisPrefix string
)

// Config use to init model package
type Config struct {
	RedisPool   *redis.Pool
	RedisPrefix string
}

// Init inits model package
func Init(config Config) error {
	redisPool = config.RedisPool
	redisPrefix = config.RedisPrefix

	return nil
}
