package internal

import (
	"time"

	"github.com/acoshift/configfile"
	"github.com/garyburd/redigo/redis"
)

var config = configfile.NewReader("config")

var (
	redisAddr = config.String("redis_addr")
	redisDB   = config.Int("redis_db")
	redisPass = config.String("redis_pass")
)

var (
	pool = &redis.Pool{
		IdleTimeout: 10 * time.Minute,
		MaxIdle:     10,
		MaxActive:   100,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr,
				redis.DialDatabase(redisDB),
				redis.DialPassword(redisPass),
			)
		},
	}
)

// GetDB returns redis connection from pool
func GetDB() redis.Conn {
	return pool.Get()
}
