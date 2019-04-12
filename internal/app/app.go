package app

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/httpmux"
	"github.com/moonrhythm/session"
	redisstore "github.com/moonrhythm/session/store/goredis"

	"github.com/acoshift/acourse/internal/app/admin"
	"github.com/acoshift/acourse/internal/app/app"
	"github.com/acoshift/acourse/internal/app/auth"
	"github.com/acoshift/acourse/internal/app/editor"
	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
)

// Handler creates new app's handler
func Handler() http.Handler {
	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := httpmux.New()
	app.Mount(m)
	auth.Mount(m)
	editor.Mount(m)
	admin.Mount(m)

	return middleware.Chain(
		sqlctx.Middleware(config.DBClient()),
		redisctx.Middleware(config.RedisClient(), config.String("redis_prefix")),
		session.Middleware(session.Config{
			Secret:   config.Bytes("session_secret"),
			Path:     "/",
			MaxAge:   7 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			SameSite: http.SameSiteLaxMode,
			Rolling:  true,
			Proxy:    true,
			Store: redisstore.New(redisstore.Config{
				Prefix: config.String("redis_prefix"),
				Client: config.RedisClient(),
			}),
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
