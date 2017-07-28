package session_test

import "time"

type mockStore struct {
	GetFunc func(string) ([]byte, error)
	SetFunc func(string, []byte, time.Duration) error
	DelFunc func(string) error
}

func (m *mockStore) Get(key string) ([]byte, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(key)
}

func (m *mockStore) Set(key string, value []byte, ttl time.Duration) error {
	if m.SetFunc == nil {
		return nil
	}
	return m.SetFunc(key, value, ttl)
}

func (m *mockStore) Del(key string) error {
	if m.DelFunc == nil {
		return nil
	}
	return m.DelFunc(key)
}
