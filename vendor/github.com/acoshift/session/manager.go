package session

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

// Manager is the session manager
type Manager struct {
	config Config
	hashID func(id string) string
}

// New creates new session manager
func New(config Config) *Manager {
	if config.Store == nil {
		panic("session: nil store")
	}

	m := Manager{}
	m.config = config

	m.hashID = func(id string) string {
		h := sha256.New()
		h.Write([]byte(id))
		h.Write(config.Secret)
		return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
	}

	if config.DisableHashID {
		m.hashID = func(id string) string {
			return id
		}
	}

	return &m
}

// Get retrieves session from request
func (m *Manager) Get(r *http.Request, name string) *Session {
	s := Session{
		DisableRenew: m.config.DisableRenew,
		Name:         name,
		Domain:       m.config.Domain,
		Path:         m.config.Path,
		HTTPOnly:     m.config.HTTPOnly,
		MaxAge:       m.config.MaxAge,
		Secure:       (m.config.Secure == ForceSecure) || (m.config.Secure == PreferSecure && isTLS(r)),
	}

	// get session key from cookie
	cookie, err := r.Cookie(name)
	if err == nil && len(cookie.Value) > 0 {
		// get session data from store
		b, err := m.config.Store.Get(m.hashID(cookie.Value))
		if err == nil {
			s.id = cookie.Value
			s.decode(b)
		}
		// DO NOT set session id to cookie value if not found in store
		// to prevent session fixation attack
	}

	return &s
}

// Save saves session to store and set cookie to response
//
// Save must be called before response header was written
func (m *Manager) Save(w http.ResponseWriter, s *Session) {
	s.setCookie(w)

	hashedID := m.hashID(s.id)
	switch s.mark.(type) {
	case markDestroy:
		m.config.Store.Del(hashedID)
	case markRotate:
		if len(s.oldID) > 0 {
			s.Set(timestampKey{}, int64(-1))
			m.config.Store.Set(m.hashID(s.oldID), s.encode(), m.config.RenewalTimeout)
		}
		s.Set(timestampKey{}, time.Now().Unix())
		m.config.Store.Set(hashedID, s.encode(), s.MaxAge)
	default:
		if s.changed {
			s.Set(timestampKey{}, time.Now().Unix())
			m.config.Store.Set(hashedID, s.encode(), s.MaxAge)
		}
	}
}
