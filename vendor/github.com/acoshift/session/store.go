package session

import (
	"time"
)

// Store interface
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Del(key string) error
	Exp(key string, ttl time.Duration) error
}
