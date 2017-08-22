package redis_test

import (
	"testing"
	"time"

	store "github.com/acoshift/session/store/redis"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	s := store.New(store.Config{Prefix: "session:", Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}})
	err := s.Set("a", []byte("test"), time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	b, err := s.Get("a")
	assert.Nil(t, b, "expected expired key return nil")
	assert.Error(t, err)

	s.Set("a", []byte("test"), time.Second)
	time.Sleep(2 * time.Second)
	_, err = s.Get("a")
	assert.Error(t, err, "expected expired key return error")

	s.Set("a", []byte("test"), time.Second)
	b, err = s.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(b))

	s.Del("a")
	_, err = s.Get("a")
	assert.Error(t, err)
}
