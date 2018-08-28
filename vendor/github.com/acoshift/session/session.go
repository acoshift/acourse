package session

import (
	"net/http"
	"time"

	"github.com/acoshift/flash"
)

// Data stores session data
type Data map[string]interface{}

// Session type
type Session struct {
	id      string // id is the hashed id if enable hash
	rawID   string
	oldID   string // for regenerate, is the hashed old id if enable hash
	oldData Data   // is the old data before regenerate
	data    Data
	destroy bool
	changed bool
	isNew   bool
	flash   *flash.Flash

	// cookie config
	Name     string
	Domain   string
	Path     string
	HTTPOnly bool
	MaxAge   time.Duration
	Secure   bool
	SameSite http.SameSite
	Rolling  bool

	manager *Manager
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

// Regenerate regenerates session id
// use when change user access level to prevent session fixation
//
// can not use regenerate and destroy same time
// Regenerate can call only one time
func (s *Session) Regenerate() {
	if len(s.oldID) > 0 {
		return
	}

	if s.destroy {
		return
	}

	s.oldID = s.id
	s.oldData = s.data.Clone()
	s.rawID = s.manager.config.GenerateID()
	s.isNew = true
	s.id = s.manager.hashID(s.rawID)
	s.changed = true
}

// IsNew checks is new session
func (s *Session) IsNew() bool {
	return s.isNew
}

// Renew clear all data in current session
// and regenerate session id
func (s *Session) Renew() {
	s.data = make(Data)
	s.Regenerate()
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

	if s.isNew && !s.Changed() {
		return
	}
	if !s.Rolling && (!s.isNew || !s.Changed()) {
		return
	}

	value := s.rawID
	if len(s.manager.config.Keys) > 0 {
		digest := sign(value, s.manager.config.Keys[0])
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

// Flash returns flash from session
func (s *Session) Flash() *flash.Flash {
	if s.flash != nil {
		return s.flash
	}
	if b, ok := s.Get(flashKey).([]byte); ok {
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
	if t, ok := s.Get(destroyedKey).(int64); ok {
		if t < time.Now().UnixNano()-int64(HijackedTime) {
			return true
		}
	}
	return false
}
