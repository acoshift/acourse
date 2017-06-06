package app

import (
	"database/sql"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	sqlstore "github.com/acoshift/session/store/sql"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/xsrftoken"
)

// Middleware wraps handlers with app's middleware
func Middleware(h http.Handler) http.Handler {
	var sessionStore session.Store
	if len(redisAddr) > 0 {
		sessionStore = redisstore.New(redisstore.Config{
			Prefix: "acourse:",
			Pool: &redis.Pool{
				MaxIdle:     20,
				IdleTimeout: 10 * time.Minute,
				Dial: func() (redis.Conn, error) {
					return redis.Dial("tcp", redisAddr, redis.DialDatabase(redisDB), redis.DialPassword(redisPass))
				},
			},
		})
	} else {
		sessionStore = sqlstore.New(sqlstore.Config{
			DB:              db,
			Table:           "sessions",
			CleanupInterval: 8 * time.Hour,
		})
	}

	return middleware.Chain(
		recovery,
		session.Middleware(session.Config{
			Name:     "sess",
			Entropy:  32,
			Path:     "/",
			MaxAge:   10 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			Store:    sessionStore,
		}),
		flash.Middleware(),
		fetchUser,
		xsrf,
	)(h)
}

func recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var p string
				switch t := r.(type) {
				case string:
					p = t
				case error:
					p = t.Error()
				default:
					p = "unknown"
				}
				debug.PrintStack()
				http.Error(w, p, http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func xsrf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id string
		if u := appctx.GetUser(r.Context()); u != nil {
			id = u.ID
		}
		if r.Method == http.MethodPost {
			x := r.FormValue("X")
			if !xsrftoken.Valid(x, xsrfSecret, id, r.URL.Path) {
				http.Error(w, "invalid xsrf token, go back, refresh and try again...", http.StatusBadRequest)
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		token := xsrftoken.Generate(xsrfSecret, id, r.URL.Path)
		ctx := appctx.WithXSRFToken(r.Context(), token)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func mustSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		id, _ := s.Get(keyUserID).(string)
		if len(id) == 0 {
			http.Redirect(w, r, "/signin?r="+url.QueryEscape(r.RequestURI), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func mustNotSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		id, _ := s.Get(keyUserID).(string)
		if len(id) > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func fetchUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := session.Get(ctx)
		id, _ := s.Get(keyUserID).(string)
		if len(id) > 0 {
			u, err := model.GetUser(id)
			if err == model.ErrNotFound {
				u = &model.User{
					ID:       id,
					Username: id,
				}
			}
			r = r.WithContext(appctx.WithUser(ctx, u))
		}
		h.ServeHTTP(w, r)
	})
}

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := appctx.GetUser(r.Context())
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Admin.Bool {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func onlyInstructor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := appctx.GetUser(r.Context())
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Instructor.Bool {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func isCourseOwner(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := appctx.GetUser(ctx)
		if u == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)

		var ownerID string
		err := db.QueryRow(`select user_id from courses where id = $1`, id).Scan(&ownerID)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ownerID != u.ID {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
