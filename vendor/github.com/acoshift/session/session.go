package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"
	"time"
)

type (
	markSave    struct{}
	markRolling struct{}
	markDestroy struct{}
	markRotate  struct{}
)

// Session type
type Session struct {
	id          string
	oldID       string // for rotate
	data        map[interface{}]interface{}
	rawData     []byte
	mark        interface{}
	encodedData []byte

	// cookie config
	entropy  int
	Name     string
	Domain   string
	Path     string
	HTTPOnly bool
	MaxAge   time.Duration
	Secure   bool
}

func init() {
	gob.Register(map[interface{}]interface{}{})
}

type sessionKey struct{}

// Get gets session from context
func Get(ctx context.Context) *Session {
	s, _ := ctx.Value(sessionKey{}).(*Session)
	return s
}

// Set sets session to context
func Set(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, s)
}

func (s *Session) encode() ([]byte, error) {
	if len(s.data) == 0 {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&s.data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Session) decode(b []byte) {
	s.data = make(map[interface{}]interface{})
	if len(b) > 0 {
		gob.NewDecoder(bytes.NewReader(b)).Decode(&s.data)
	}
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
	s.data[key] = value
}

// Del deletes data from session
func (s *Session) Del(key interface{}) {
	if s.data == nil {
		return
	}
	delete(s.data, key)
}

// Rotate rotates session id
// use when change user access level to prevent session fixation
//
// can not use rotate and destory same time
func (s *Session) Rotate() {
	s.mark = markRotate{}
}

// Destroy destroys session from store
func (s *Session) Destroy() {
	s.mark = markDestroy{}
}

func (s *Session) setCookie(w http.ResponseWriter) {
	if _, ok := s.mark.(markDestroy); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     s.Name,
			Domain:   s.Domain,
			Path:     s.Path,
			HttpOnly: s.HTTPOnly,
			Value:    "",
			MaxAge:   -1,
			Secure:   s.Secure,
		})
		return
	}

	// set cookie only if session value changed
	var err error
	s.encodedData, err = s.encode()
	if err != nil {
		// this should never happended
		// or developer don't register struct into gob
		return
	}
	if _, ok := s.mark.(markRotate); ok {
		s.oldID = s.id
		s.id = generateID(s.entropy)
	} else if bytes.Compare(s.rawData, s.encodedData) == 0 {
		if len(s.encodedData) == 0 {
			// empty session
			return
		}
		// should rolling cookie
		s.mark = markRolling{}
	} else {
		s.mark = markSave{}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.Name,
		Domain:   s.Domain,
		Path:     s.Path,
		HttpOnly: s.HTTPOnly,
		Value:    s.id,
		MaxAge:   int(s.MaxAge / time.Second),
		Secure:   s.Secure,
	})
}
