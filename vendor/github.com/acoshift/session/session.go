package session

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/acoshift/flash"
)

// Session type
type Session struct {
	id      string // id is the hashed id if enable hash
	rawID   string
	oldID   string // for rotate, is the hashed old id if enable hash
	oldData []byte // is the old encoded data before rotate
	data    map[interface{}]interface{}
	destroy bool
	changed bool
	flash   *flash.Flash

	// cookie config
	Name     string
	Domain   string
	Path     string
	HTTPOnly bool
	MaxAge   time.Duration
	Secure   bool
	SameSite SameSite

	IDHashFunc func(id string) string
}

func init() {
	gob.Register(map[interface{}]interface{}{})
	gob.Register(flashKey{})
}

// session internal data
type (
	flashKey struct{}
)

func (s *Session) encode() []byte {
	if len(s.data) == 0 {
		return []byte{}
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(s.data)
	if err != nil {
		// this should never happened
		// or developer don't register type into gob
		panic("session: can not encode data; " + err.Error())
	}
	return buf.Bytes()
}

func (s *Session) decode(b []byte) {
	s.data = make(map[interface{}]interface{})
	if len(b) > 0 {
		gob.NewDecoder(bytes.NewReader(b)).Decode(&s.data)
	}
}

// ID returns session id or hashed session id if enable hash id
func (s *Session) ID() string {
	return s.id
}

// Changed returns is session data changed
func (s *Session) Changed() bool {
	if s.changed {
		return true
	}
	if s.flash != nil && s.flash.Changed() {
		s.changed = true
		return true
	}
	return false
}

// Get gets data from session
func (s *Session) Get(key interface{}) interface{} {
	if s.data == nil {
		return nil
	}
	return s.data[key]
}

// Set sets data to session
func (s *Session) Set(key, value interface{}) {
	if s.data == nil {
		s.data = make(map[interface{}]interface{})
	}
	s.changed = true
	s.data[key] = value
}

// Del deletes data from session
func (s *Session) Del(key interface{}) {
	if s.data == nil {
		return
	}
	if _, ok := s.data[key]; ok {
		s.changed = true
		delete(s.data, key)
	}
}

// Pop gets data from session then delete it
func (s *Session) Pop(key interface{}) interface{} {
	if s.data == nil {
		return nil
	}
	r := s.data[key]
	s.changed = true
	delete(s.data, key)
	return r
}

// Rotate rotates session id
// use when change user access level to prevent session fixation
//
// can not use rotate and destroy same time
// Rotate can call only one time
func (s *Session) Rotate() {
	if len(s.oldID) > 0 {
		return
	}

	if s.destroy {
		return
	}

	s.oldID = s.id
	s.oldData = s.encode()
	s.rawID = generateID()
	if s.IDHashFunc != nil {
		s.id = s.IDHashFunc(s.rawID)
	} else {
		s.id = s.rawID
	}
	s.changed = true
}

// Renew clear all data in current session
// and rotate session id
func (s *Session) Renew() {
	s.data = make(map[interface{}]interface{})
	s.Rotate()
}

// Destroy destroys session from store
func (s *Session) Destroy() {
	s.destroy = true
}

func (s *Session) setCookie(w http.ResponseWriter) {
	if s.destroy {
		http.SetCookie(w, &http.Cookie{
			Name:     s.Name,
			Domain:   s.Domain,
			Path:     s.Path,
			HttpOnly: s.HTTPOnly,
			Value:    "",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			Secure:   s.Secure,
		})
		return
	}

	// if session don't have raw id, don't set cookie
	if len(s.rawID) == 0 {
		return
	}

	// if session not modified, don't set cookie
	if !s.Changed() {
		return
	}

	setCookie(w, &cookie{
		Cookie: http.Cookie{
			Name:     s.Name,
			Domain:   s.Domain,
			Path:     s.Path,
			HttpOnly: s.HTTPOnly,
			Value:    s.rawID,
			MaxAge:   int(s.MaxAge / time.Second),
			Expires:  time.Now().Add(s.MaxAge),
			Secure:   s.Secure,
		},
		SameSite: s.SameSite,
	})
}

// Flash returns flash from session
func (s *Session) Flash() *flash.Flash {
	if s.flash != nil {
		return s.flash
	}
	if b, ok := s.Get(flashKey{}).([]byte); ok {
		s.flash, _ = flash.Decode(b)
	}
	if s.flash == nil {
		s.flash = flash.New()
	}
	return s.flash
}

// Hijacked checks is session was hijacked,
// can use only with Manager
func (s *Session) Hijacked() bool {
	if t, ok := s.Get(destroyedKey{}).(int64); ok {
		if t < time.Now().UnixNano()-int64(HijackedTime) {
			return true
		}
	}
	return false
}
