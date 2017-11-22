package session

import (
	"time"
)

// Config is the session manager config
type Config struct {
	Store  Store
	Secret []byte // session id salt when put to store

	// Cookie config
	Domain   string
	HTTPOnly bool
	Path     string
	MaxAge   time.Duration
	Secure   Secure
	SameSite SameSite

	// DeleteOldSession deletes the old session from store when rotate,
	// better not to delete old session to avoid user loss session when unstable network
	DeleteOldSession bool

	// Disable features
	DisableRenew  bool // disable auto renew session
	DisableHashID bool // disable hash session id when save to store
}

// Secure config
type Secure int

// Secure values
const (
	NoSecure     Secure = iota
	PreferSecure        // if request is https will set secure cookie
	ForceSecure         // always set secure cookie
)

// SameSite config
//
// http://httpwg.org/http-extensions/draft-ietf-httpbis-cookie-same-site.html
type SameSite string

// SameSite values
const (
	SameSiteNone   SameSite = ""
	SameSiteLax    SameSite = "Lax"
	SameSiteStrict SameSite = "Strict"
)

// Global Session Config
var (
	HijackedTime = 5 * time.Minute
)
