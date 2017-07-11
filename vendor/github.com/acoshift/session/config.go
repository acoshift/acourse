package session

import (
	"time"
)

// Config is the session middleware config
type Config struct {
	Store   Store
	Entropy int // session id entropy in byte, default is 32 byte

	// Cookie config
	Name     string // Cookie name, default is "sess"
	Domain   string
	HTTPOnly bool
	Path     string
	MaxAge   time.Duration
	Secure   Secure
}

// Secure config
type Secure int

// Secure configs
const (
	NoSecure     Secure = iota
	PreferSecure        // if request is https will set secure cookie
	ForceSecure         // always set secure cookie
)
