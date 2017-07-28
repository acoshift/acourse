package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/acoshift/session"
)

// Config is the memory store config
type Config struct {
	CleanupInterval time.Duration
}

// New creates new memory store
func New(config Config) session.Store {
	s := &memoryStore{
		cleanupInterval: config.CleanupInterval,
		l:               make(map[interface{}]*item),
	}
	if config.CleanupInterval > 0 {
		go s.cleanupWorker()
	}
	return s
}

type item struct {
	data []byte
	exp  time.Time
}

type memoryStore struct {
	cleanupInterval time.Duration
	m               sync.RWMutex
	l               map[interface{}]*item
}

func (s *memoryStore) cleanupWorker() {
	now := time.Now()
	s.m.Lock()
	for k, v := range s.l {
		if !v.exp.IsZero() && v.exp.Before(now) {
			delete(s.l, k)
		}
	}
	s.m.Unlock()
	time.AfterFunc(s.cleanupInterval, s.cleanupWorker)
}

var errNotFound = errors.New("memory: session not found")

func (s *memoryStore) Get(key string) ([]byte, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	v := s.l[key]
	if v == nil {
		return nil, errNotFound
	}
	if !v.exp.IsZero() && v.exp.Before(time.Now()) {
		return nil, errNotFound
	}
	return v.data, nil
}

func (s *memoryStore) Set(key string, value []byte, ttl time.Duration) error {
	s.m.Lock()
	s.l[key] = &item{
		data: value,
		exp:  time.Now().Add(ttl),
	}
	s.m.Unlock()
	return nil
}

func (s *memoryStore) Del(key string) error {
	s.m.Lock()
	delete(s.l, key)
	s.m.Unlock()
	return nil
}
