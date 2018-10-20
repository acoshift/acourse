package session

import (
	"net/http"
	"time"
)

// Data stores session data
type Data map[string]interface{}

// Session type
type Session struct {
	id      string // id is the hashed id if hash enabled
	rawID   string
	data    Data
	changed bool
	isNew   bool
	flash   *Flash

	// cookie config
	Name     string
	Domain   string
	Path     string
	HTTPOnly bool
	MaxAge   time.Duration
	Secure   bool
	SameSite http.SameSite
	Rolling  bool

	m *scopedManager
}

// Clone clones session data
func (data Data) Clone() Data {
	r := make(Data)
	for k, v := range data {
		r[k] = v
	}
	return r
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
func (s *Session) Get(key string) interface{} {
	if s.data == nil {
		return nil
	}
	return s.data[key]
}

// GetString gets string from session
func (s *Session) GetString(key string) string {
	r, _ := s.Get(key).(string)
	return r
}

// GetInt gets int from session
func (s *Session) GetInt(key string) int {
	r, _ := s.Get(key).(int)
	return r
}

// GetInt64 gets int64 from session
func (s *Session) GetInt64(key string) int64 {
	r, _ := s.Get(key).(int64)
	return r
}

// GetFloat32 gets float32 from session
func (s *Session) GetFloat32(key string) float32 {
	r, _ := s.Get(key).(float32)
	return r
}

// GetFloat64 gets float64 from session
func (s *Session) GetFloat64(key string) float64 {
	r, _ := s.Get(key).(float64)
	return r
}

// GetBool gets bool from session
func (s *Session) GetBool(key string) bool {
	r, _ := s.Get(key).(bool)
	return r
}

// Set sets data to session
func (s *Session) Set(key string, value interface{}) {
	if s.data == nil {
		s.data = make(Data)
	}
	s.changed = true
	s.data[key] = value
}

// Del deletes data from session
func (s *Session) Del(key string) {
	if s.data == nil {
		return
	}
	if _, ok := s.data[key]; ok {
		s.changed = true
		delete(s.data, key)
	}
}

// Pop gets data from session then delete it
func (s *Session) Pop(key string) interface{} {
	if s.data == nil {
		return nil
	}

	r, ok := s.data[key]
	if ok {
		s.changed = true
		delete(s.data, key)
	}
	return r
}

// PopString pops string from session
func (s *Session) PopString(key string) string {
	r, _ := s.Pop(key).(string)
	return r
}

// PopInt pops int from session
func (s *Session) PopInt(key string) int {
	r, _ := s.Pop(key).(int)
	return r
}

// PopInt64 pops int64 from session
func (s *Session) PopInt64(key string) int64 {
	r, _ := s.Pop(key).(int64)
	return r
}

// PopFloat32 pops float32 from session
func (s *Session) PopFloat32(key string) float32 {
	r, _ := s.Pop(key).(float32)
	return r
}

// PopFloat64 pops float64 from session
func (s *Session) PopFloat64(key string) float64 {
	r, _ := s.Pop(key).(float64)
	return r
}

// PopBool pops bool from session
func (s *Session) PopBool(key string) bool {
	r, _ := s.Pop(key).(bool)
	return r
}

// IsNew checks is new session
func (s *Session) IsNew() bool {
	return s.isNew
}

// Flash returns flash from session,
func (s *Session) Flash() *Flash {
	if s.flash != nil {
		return s.flash
	}

	s.flash = new(Flash)
	if b, ok := s.Get(flashKey).([]byte); ok {
		s.flash.decode(b)
	}
	return s.flash
}

// Hijacked checks is session was hijacked
func (s *Session) Hijacked() bool {
	if t, ok := s.Get(destroyedKey).(int64); ok {
		if t < time.Now().UnixNano()-int64(HijackedTime) {
			return true
		}
	}
	return false
}

// with scopedManager

// Regenerate regenerates session id
// use when change user access level to prevent session fixation
//
// Can use only with middleware
func (s *Session) Regenerate() error {
	if s.m == nil {
		return ErrNotPassMiddleware
	}
	return s.m.Regenerate(s)
}

// Renew clear all data in current session and regenerate session id
//
// Can use only with middleware
func (s *Session) Renew() error {
	if s.m == nil {
		return ErrNotPassMiddleware
	}
	return s.m.Renew(s)
}

// Destroy destroys session from store
//
// Can use only with middleware
func (s *Session) Destroy() error {
	if s.m == nil {
		return ErrNotPassMiddleware
	}
	return s.m.Destroy(s)
}
