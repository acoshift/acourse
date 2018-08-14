package session

import (
	"errors"
	"time"
)

// Errors
var (
	// ErrNotFound is the error when session not found
	// store must return ErrNotFound if session data not exists
	ErrNotFound = errors.New("session: not found")
)

// Store interface
type Store interface {
	Get(key string, opt StoreOption) (Data, error)
	Set(key string, value Data, opt StoreOption) error
	Del(key string, opt StoreOption) error
}

// StoreOption type
type StoreOption struct {
	Rolling bool
	TTL     time.Duration
}

func makeStoreOption(m *Manager, s *Session) StoreOption {
	return StoreOption{
		Rolling: s.Rolling,
		TTL:     m.config.IdleTimeout,
	}
}
