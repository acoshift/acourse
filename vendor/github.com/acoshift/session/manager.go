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

	if config.DisableHashID {
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

	return &m
}

// Get retrieves session from request
func (m *Manager) Get(r *http.Request, name string) *Session {
	s := Session{
		manager:  m,
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

	return &s
}

// Save saves session to store and set cookie to response
//
// Save must be called before response header was written
func (m *Manager) Save(w http.ResponseWriter, s *Session) error {
	s.setCookie(w)

	opt := makeStoreOption(m, s)

	if s.destroy {
		return m.config.Store.Del(s.id, opt)
	}

	// detect is flash changed and encode new flash data
	if s.flash != nil && s.flash.Changed() {
		b, _ := s.flash.Encode()
		s.Set(flashKey, b)
	}

	// if session not modified, don't save to store to prevent store overflow
	if !m.config.Resave && !s.Changed() {
		return nil
	}

	// check is regenerate
	if len(s.oldID) > 0 {
		if m.config.DeleteOldSession {
			err := m.config.Store.Del(s.oldID, opt)
			if err != nil {
				return err
			}
		} else {
			// save old session data if not delete
			s.oldData[timestampKey] = int64(0)
			s.oldData[destroyedKey] = time.Now().UnixNano()
			err := m.config.Store.Set(s.oldID, s.oldData, opt)
			if err != nil {
				return err
			}
		}
	}

	// save sesion data to store
	s.Set(timestampKey, time.Now().Unix())
	return m.config.Store.Set(s.id, s.data, opt)
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
