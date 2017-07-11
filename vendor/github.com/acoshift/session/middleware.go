package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/acoshift/middleware"
)

// Middleware is the session parser middleware
func Middleware(config Config) middleware.Middleware {
	if config.Store == nil {
		panic("session: nil store")
	}

	// set default config
	if config.Entropy <= 0 {
		config.Entropy = 32
	}

	if len(config.Name) == 0 {
		config.Name = "sess"
	}

	generateID := func() string {
		b := make([]byte, config.Entropy)
		if _, err := rand.Read(b); err != nil {
			// this should never happended
			// or something wrong with OS's crypto pseudorandom generator
			panic(err)
		}
		return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
	}

	hashID := func(id string) string {
		h := sha256.New()
		h.Write([]byte(id))
		h.Write(config.Secret)
		return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := Session{
				generateID: generateID,
				Name:       config.Name,
				Domain:     config.Domain,
				Path:       config.Path,
				HTTPOnly:   config.HTTPOnly,
				MaxAge:     config.MaxAge,
				Secure:     (config.Secure == ForceSecure) || (config.Secure == PreferSecure && isTLS(r)),
			}

			// get session key from cookie
			cookie, err := r.Cookie(config.Name)
			if err == nil && len(cookie.Value) > 0 {
				// get session data from store
				s.rawData, err = config.Store.Get(hashID(cookie.Value))
				if err == nil {
					s.id = cookie.Value
					s.decode(s.rawData)
				}
				// DO NOT set session id to cookie value if not found in store
				// to prevent session fixation attack
			}

			// if session not found, create new session
			if len(s.id) == 0 {
				s.id = generateID()
			}

			// use defer to alway save session even panic
			defer func() {
				hID := hashID(s.id)
				switch s.mark.(type) {
				case markDestroy:
					config.Store.Del(hID)
				case markSave:
					// if session was modified, save session to store,
					// if not don't save to store to prevent store overflow
					config.Store.Set(hID, s.encodedData, s.MaxAge)
				case markRolling:
					// session not modified but not empty
					config.Store.Exp(hID, config.MaxAge)
				case markRotate:
					config.Store.Set(hID, s.encodedData, s.MaxAge)
					config.Store.Del(hashID(s.oldID))
				}
			}()

			nr := r.WithContext(Set(r.Context(), &s))
			nw := sessionWriter{
				ResponseWriter: w,
				s:              &s,
			}
			h.ServeHTTP(&nw, nr)
		})
	}
}
