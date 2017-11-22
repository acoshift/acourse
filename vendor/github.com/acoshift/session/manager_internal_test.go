package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldRenew(t *testing.T) {
	m := Manager{
		config: Config{
			DisableRenew: false,
		},
	}

	s := &Session{}
	s.Set(timestampKey{}, int64(-1))
	assert.False(t, m.shouldRenewSession(s), "expected sec -1 should not renew")

	s.Set(timestampKey{}, int64(0))
	assert.False(t, m.shouldRenewSession(s), "expected sec 0 should not renew")

	now := time.Now().Unix()

	s.MaxAge = 10 * time.Second
	s.Set(timestampKey{}, now-7)
	assert.True(t, m.shouldRenewSession(s), "expected sec -7 of max-age 10 should renew")

	s.Set(timestampKey{}, now-5)
	assert.True(t, m.shouldRenewSession(s), "expected sec -5 of max-age 10 should renew")

	s.Set(timestampKey{}, now-3)
	assert.False(t, m.shouldRenewSession(s), "expected sec -3 of max-age 10 should not renew")

	// test disable renew
	m.config.DisableRenew = true
	s.Set(timestampKey{}, now-7)
	assert.False(t, m.shouldRenewSession(s), "expected disable renew should not renew")
}
