package session

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"strings"
	"time"
)

// Manager is the session manager
type Manager struct {
	config Config
	hashID func(id string) string
}

// manager internal data
type (
	timestampKey struct{}
	destroyedKey struct{} // for detect session hijack
)

func init() {
	gob.Register(timestampKey{})
	gob.Register(destroyedKey{})
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
		Name:       name,
		Domain:     m.config.Domain,
		Path:       m.config.Path,
		HTTPOnly:   m.config.HTTPOnly,
		MaxAge:     m.config.MaxAge,
		Secure:     (m.config.Secure == ForceSecure) || (m.config.Secure == PreferSecure && isTLS(r)),
		SameSite:   m.config.SameSite,
		IDHashFunc: m.hashID,
	}

	// get session key from cookie
	cookie, err := r.Cookie(name)
	if err == nil && len(cookie.Value) > 0 {
		hashedID := m.hashID(cookie.Value)

		// get session data from store
		b, err := m.config.Store.Get(hashedID)
		if err == nil {
			s.id = hashedID
			s.decode(b)
		}
		// DO NOT set session id to cookie value if not found in store
		// to prevent session fixation attack
	}

	if len(s.id) == 0 {
		s.rawID = generateID()
		s.id = m.hashID(s.rawID)
	}

	return &s
}

// Save saves session to store and set cookie to response
//
// Save must be called before response header was written
func (m *Manager) Save(w http.ResponseWriter, s *Session) {
	// check is session should renew
	if m.shouldRenewSession(s) {
		// use rotate to renew session
		s.Rotate()
	}

	s.setCookie(w)

	if s.destroy {
		m.config.Store.Del(s.id)
		return
	}

	// detect is flash changed and encode new flash data
	if s.flash != nil && s.flash.Changed() {
		b, _ := s.flash.Encode()
		s.Set(flashKey{}, b)
	}

	// if session not modified, don't save to store to prevent store overflow
	if !s.Changed() {
		return
	}

	// check is rotate
	if len(s.oldID) > 0 {
		if m.config.DeleteOldSession {
			m.config.Store.Del(s.oldID)
		} else {
			// save old session data if not delete
			var d Session
			d.decode(s.oldData)
			d.Set(timestampKey{}, int64(0))
			d.Set(destroyedKey{}, time.Now().UnixNano())
			m.config.Store.Set(s.oldID, d.encode(), s.MaxAge)
		}
	}

	// save sesion data to store
	s.Set(timestampKey{}, time.Now().Unix())
	m.config.Store.Set(s.id, s.encode(), s.MaxAge)
}

func (m *Manager) shouldRenewSession(s *Session) bool {
	if m.config.DisableRenew {
		return false
	}
	sec, _ := s.Get(timestampKey{}).(int64)
	if sec <= 0 {
		return false
	}
	t := time.Unix(sec, 0)
	if time.Now().Sub(t) < s.MaxAge/2 {
		return false
	}
	return true
}
