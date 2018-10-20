package session

import (
	"crypto/rand"
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

	if m.config.GenerateID == nil {
		m.config.GenerateID = func() string {
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				// this should never happened
				// or something wrong with OS's crypto pseudorandom generator
				panic(err)
			}
			return base64.RawURLEncoding.EncodeToString(b)
		}
	}

	if m.config.DisableHashID {
		m.hashID = func(id string) string {
			return id
		}
	} else {
		m.hashID = func(id string) string {
			h := sha256.New()
			h.Write([]byte(id))
			h.Write(config.Secret)
			return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
		}
	}

	if m.config.IdleTimeout <= 0 {
		m.config.IdleTimeout = m.config.MaxAge
	}

	return &m
}

// Get retrieves session from request
func (m *Manager) Get(r *http.Request, name string) (*Session, error) {
	s := Session{
		Name:     name,
		Domain:   m.config.Domain,
		Path:     m.config.Path,
		HTTPOnly: m.config.HTTPOnly,
		MaxAge:   m.config.MaxAge,
		Secure:   m.isSecure(r),
		SameSite: m.config.SameSite,
		Rolling:  m.config.Rolling,
	}

	// get session id from cookie
	cookie, err := r.Cookie(name)
	if err == nil && len(cookie.Value) > 0 {
		var rawID string

		// verify signature
		if len(m.config.Keys) > 0 {
			parts := strings.Split(cookie.Value, ".")
			rawID = parts[0]

			if len(parts) != 2 || !verify(rawID, parts[1], m.config.Keys) {
				goto invalidSignature
			}
		} else {
			rawID = cookie.Value
		}

		hashedID := m.hashID(rawID)

		// get session data from store
		s.data, err = m.config.Store.Get(hashedID, makeStoreOption(m, &s))
		if err == nil {
			s.rawID = rawID
			s.id = hashedID
		} else if err != ErrNotFound {
			return nil, err
		}

		// DO NOT set session id to cookie value if not found in store
		// to prevent session fixation attack
	}
invalidSignature:

	if len(s.id) == 0 {
		s.rawID = m.config.GenerateID()
		s.id = m.hashID(s.rawID)
		s.isNew = true
	}

	return &s, nil
}

// Save saves session to store and set cookie to response
//
// Save must be called before response header was written
func (m *Manager) Save(w http.ResponseWriter, s *Session) error {
	m.setCookie(w, s)

	opt := makeStoreOption(m, s)

	// detect is flash changed and encode new flash data
	if s.flash != nil && s.flash.Changed() {
		b, _ := s.flash.encode()
		s.Set(flashKey, b)
	}

	// if session not modified, don't save to store to prevent store overflow
	if !m.config.Resave && !s.Changed() {
		return nil
	}

	// save sesion data to store
	s.Set(timestampKey, time.Now().Unix())
	return m.config.Store.Set(s.id, s.data, opt)
}

// Destroy deletes session from store
func (m *Manager) Destroy(s *Session) error {
	return m.config.Store.Del(s.id, makeStoreOption(m, s))
}

// Regenerate regenerates session id
// use when change user access level to prevent session fixation
func (m *Manager) Regenerate(w http.ResponseWriter, s *Session) error {
	opt := makeStoreOption(m, s)

	id := s.id

	s.rawID = m.config.GenerateID()
	s.isNew = true
	s.id = m.hashID(s.rawID)
	s.changed = true

	if m.config.DeleteOldSession {
		return m.config.Store.Del(id, opt)
	}

	data := s.data.Clone()
	data[timestampKey] = int64(0)
	data[destroyedKey] = time.Now().UnixNano()
	return m.config.Store.Set(id, data, opt)
}

// Renew clears session data and regenerate new session id
func (m *Manager) Renew(w http.ResponseWriter, s *Session) error {
	s.data = make(Data)
	return m.Regenerate(w, s)
}

func (m *Manager) setCookie(w http.ResponseWriter, s *Session) {
	// if session don't have raw id, don't set cookie
	if len(s.rawID) == 0 {
		return
	}

	if s.isNew && !s.Changed() {
		return
	}
	if !s.Rolling && (!s.isNew || !s.Changed()) {
		return
	}

	value := s.rawID
	if len(m.config.Keys) > 0 {
		digest := sign(value, m.config.Keys[0])
		value += "." + digest
	}

	cs := http.Cookie{
		Name:     s.Name,
		Domain:   s.Domain,
		Path:     s.Path,
		HttpOnly: s.HTTPOnly,
		Value:    value,
		Secure:   s.Secure,
		SameSite: s.SameSite,
	}
	if s.MaxAge > 0 {
		cs.MaxAge = int(s.MaxAge / time.Second)
		cs.Expires = time.Now().Add(s.MaxAge)
	}

	http.SetCookie(w, &cs)
}

func (m *Manager) isSecure(r *http.Request) bool {
	if m.config.Secure == ForceSecure {
		return true
	}
	if m.config.Secure == PreferSecure {
		if r.TLS != nil {
			return true
		}
		if m.config.Proxy && r.Header.Get("X-Forwarded-Proto") == "https" {
			return true
		}
	}

	return false
}
