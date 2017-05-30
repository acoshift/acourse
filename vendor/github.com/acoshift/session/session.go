package session

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/acoshift/middleware"
)

type (
	sessionKey struct{}
)

// Middleware is the session parser middleware
func Middleware(config Config) middleware.Middleware {
	if config.Store == nil {
		panic("session: nil store")
	}

	entropy := config.Entropy
	if entropy <= 0 {
		entropy = 16
	}

	sessName := config.Name
	if len(sessName) == 0 {
		sessName = "sess"
	}

	maxAge := int(config.MaxAge / time.Second)

	generateID := func() string {
		b := make([]byte, entropy)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sess Session

			// get session key from cookie
			cookie, err := r.Cookie(sessName)
			if err == nil && len(cookie.Value) > 0 {
				// get session data from store
				sessData, err := config.Store.Get(cookie.Value)
				if err == nil && len(sessData) > 0 {
					sess.id = cookie.Value
					sess.decode(sessData)
				}
			}

			// if session not found, create new session
			if len(sess.id) == 0 {
				sess.id = generateID()
			}

			secure := false
			if (config.Secure == ForceSecure) || (config.Secure == PreferSecure && isTLS(r)) {
				secure = true
			}

			// rolling session
			http.SetCookie(w, &http.Cookie{
				Name:     sessName,
				Domain:   config.Domain,
				Path:     config.Path,
				HttpOnly: config.HTTPOnly,
				Value:    sess.id,
				MaxAge:   maxAge,
				Secure:   secure,
			})

			defer func() {
				// if session was modified, save session to store,
				// if not don't save to store to prevent brute force attack
				b, err := sess.encode()
				if err == nil {
					if bytes.Compare(sess.p, b) == 0 {
						config.Store.Exp(sess.id, config.MaxAge)
						return
					}
					config.Store.Set(sess.id, b, config.MaxAge)
				}
			}()

			nr := r.WithContext(Set(r.Context(), &sess))
			h.ServeHTTP(w, nr)
		})
	}
}

// Session type
type Session struct {
	id string
	d  map[interface{}]interface{}
	p  []byte
}

func init() {
	gob.Register(map[interface{}]interface{}{})
}

// Get gets session from context
func Get(ctx context.Context) *Session {
	sess, _ := ctx.Value(sessionKey{}).(*Session)
	return sess
}

// Set sets session to context
func Set(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, s)
}

func (sess *Session) encode() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&sess.d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (sess *Session) decode(b []byte) {
	sess.d = make(map[interface{}]interface{})
	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&sess.d)
	if err != nil {
		return
	}
}

// Get gets data from session
func (sess *Session) Get(key interface{}) interface{} {
	if sess.d == nil {
		return nil
	}
	return sess.d[key]
}

// Set sets data to session
func (sess *Session) Set(key, value interface{}) {
	if sess.d == nil {
		sess.d = make(map[interface{}]interface{})
	}
	sess.d[key] = value
}

// Del deletes data from session
func (sess *Session) Del(key interface{}) {
	if sess.d == nil {
		return
	}
	delete(sess.d, key)
}
