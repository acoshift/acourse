package session

import (
	"net/http"
	"time"
)

// Config is the session manager config
type Config struct {
	Store Store

	// Secret is the salt for hash session id before put to store
	Secret []byte

	// Keys is the keys to sign session id
	Keys [][]byte

	// Cookie config
	Domain   string
	HTTPOnly bool
	Path     string
	MaxAge   time.Duration
	Secure   Secure
	SameSite http.SameSite

	// IdleTimeout is the ttl for storage
	// if IdleTimeout is zero, it will use MaxAge
	IdleTimeout time.Duration

	// DeleteOldSession deletes the old session from store when regenerate,
	// better not to delete old session to avoid user loss session when unstable network
	DeleteOldSession bool

	// Resave forces session to save to store even if session was not modified
	Resave bool

	// Rolling, set cookie every responses
	Rolling bool

	// Proxy, also checks X-Forwarded-Proto when use prefer secure
	Proxy bool

	// DisablaHashID disables hash session id when save to store
	DisableHashID bool

	// GenerateID is the generate id function
	GenerateID func() string
}

// Secure config
type Secure int

// Secure values
const (
	NoSecure     Secure = iota
	PreferSecure        // if request is https will set secure cookie
	ForceSecure         // always set secure cookie
)

// Global Session Config
var (
	HijackedTime = 5 * time.Minute
)
