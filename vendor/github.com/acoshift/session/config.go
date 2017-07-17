package session

import (
	"time"
)

// Config is the session middleware config
type Config struct {
	Store  Store
	Secret []byte // session id salt when put to store

	// Cookie config
	Name     string // Cookie name, default is "sess"
	Domain   string
	HTTPOnly bool
	Path     string
	MaxAge   time.Duration
	Secure   Secure

	// Disable features
	DisableRenew bool // disable auto renew session
}

// Secure config
type Secure int

// Secure configs
const (
	NoSecure     Secure = iota
	PreferSecure        // if request is https will set secure cookie
	ForceSecure         // always set secure cookie
)
