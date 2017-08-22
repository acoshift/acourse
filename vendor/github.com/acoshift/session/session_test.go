package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncodeEmpty(t *testing.T) {
	s := Session{}
	b := s.encode()
	assert.NotNil(t, b, "expected encode always return not nil")
	assert.Len(t, b, 0)
}

func TestEncodeUnregisterType(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err, "expected encode unregister type panic")
	}()
	type a struct{}
	s := Session{}
	s.Set("a", a{})
	s.encode()
}

func TestSessionOperation(t *testing.T) {
	s := Session{}
	assert.Nil(t, s.Get("a"), "expected get data from empty session return nil")

	s.Del("a")
	assert.Nil(t, s.data)

	s.Set("a", 1)
	assert.Equal(t, 1, s.Get("a"))

	s.Del("a")
	assert.Nil(t, s.Get("a"), "expected get data after delete to be nil")
}

func TestShouldRenew(t *testing.T) {
	s := Session{}
	s.Set(timestampKey{}, int64(-1))
	assert.False(t, s.shouldRenew(), "expected sec -1 should not renew")

	s.Set(timestampKey{}, int64(0))
	assert.True(t, s.shouldRenew(), "expected sec 0 should renew")

	now := time.Now().Unix()

	s.MaxAge = 10 * time.Second
	s.Set(timestampKey{}, now-7)
	assert.True(t, s.shouldRenew(), "expected sec -7 of max-age 10 should renew")

	s.Set(timestampKey{}, now-5)
	assert.True(t, s.shouldRenew(), "expected sec -5 of max-age 10 should renew")

	s.Set(timestampKey{}, now-3)
	assert.False(t, s.shouldRenew(), "expected sec -3 of max-age 10 should not renew")
}

func TestRenew(t *testing.T) {
	s := Session{}
	s.Set("a", "1")
	s.Renew()
	assert.Nil(t, s.Get("a"), "expected renew must delete all data")
}
