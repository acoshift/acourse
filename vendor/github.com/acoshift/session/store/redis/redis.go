package redis

import (
	"time"

	"github.com/acoshift/session"
	"github.com/garyburd/redigo/redis"
)

// Config is the redis store config
type Config struct {
	Pool   *redis.Pool
	Prefix string
}

// New creates new redis store
func New(config Config) session.Store {
	return &redisStore{
		pool:   config.Pool,
		prefix: config.Prefix,
	}
}

type redisStore struct {
	pool   *redis.Pool
	prefix string
}

func (s *redisStore) Get(key string) ([]byte, error) {
	c := s.pool.Get()
	defer c.Close()
	return redis.Bytes(c.Do("GET", s.prefix+key))
}

func (s *redisStore) Set(key string, value []byte, ttl time.Duration) error {
	c := s.pool.Get()
	defer c.Close()
	var err error
	if ttl > 0 {
		_, err = c.Do("SETEX", s.prefix+key, int64(ttl/time.Second), value)
	} else {
		_, err = c.Do("SET", s.prefix+key, value)
	}
	return err
}

func (s *redisStore) Del(key string) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", s.prefix+key)
	return err
}
