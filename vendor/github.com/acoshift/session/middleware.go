package session

import (
	"net/http"

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

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := Session{
				entropy:  config.Entropy,
				Name:     config.Name,
				Domain:   config.Domain,
				Path:     config.Path,
				HTTPOnly: config.HTTPOnly,
				MaxAge:   config.MaxAge,
				Secure:   (config.Secure == ForceSecure) || (config.Secure == PreferSecure && isTLS(r)),
			}

			// get session key from cookie
			cookie, err := r.Cookie(config.Name)
			if err == nil && len(cookie.Value) > 0 {
				// get session data from store
				s.rawData, err = config.Store.Get(cookie.Value)
				if err == nil {
					s.id = cookie.Value
					s.decode(s.rawData)
				}
				// DO NOT set session id to cookie value if not found in store
				// to prevent session fixation attack
			}

			// if session not found, create new session
			if len(s.id) == 0 {
				s.id = generateID(s.entropy)
			}

			// use defer to alway save session even panic
			defer func() {
				switch s.mark.(type) {
				case markDestroy:
					config.Store.Del(s.id)
				case markSave:
					// if session was modified, save session to store,
					// if not don't save to store to prevent store overflow
					config.Store.Set(s.id, s.encodedData, s.MaxAge)
				case markRolling:
					// session not modified but not empty
					config.Store.Exp(s.id, config.MaxAge)
				case markRotate:
					config.Store.Set(s.id, s.encodedData, s.MaxAge)
					config.Store.Del(s.oldID)
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
