package app

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/go-redis/redis"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/session"

	"github.com/acoshift/acourse/internal/app/admin"
	"github.com/acoshift/acourse/internal/app/app"
	"github.com/acoshift/acourse/internal/app/auth"
	"github.com/acoshift/acourse/internal/app/editor"
	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
)

// Config is the app's config
type Config struct {
	DB            *sql.DB
	SessionSecret []byte
	RedisClient   *redis.Client
	RedisPrefix   string
	SessionStore  session.Store
}

// Handler creates new app's handler
func Handler(c Config) http.Handler {
	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := http.NewServeMux()
	app.Mount(m)
	auth.Mount(m)
	editor.Mount(m)
	admin.Mount(m)

	return middleware.Chain(
		sqlctx.Middleware(c.DB),
		redisctx.Middleware(c.RedisClient, c.RedisPrefix),
		session.Middleware(session.Config{
			Secret:   c.SessionSecret,
			Path:     "/",
			MaxAge:   7 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			SameSite: http.SameSiteLaxMode,
			Rolling:  true,
			Proxy:    true,
			Store:    c.SessionStore,
		}),
		turbolinks,
		appctx.Middleware(),
	)(m)
}

func turbolinks(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Turbolinks-Location", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
