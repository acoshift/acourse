package app

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"cloud.google.com/go/errorreporting"
	"github.com/acoshift/configfile"
	"github.com/acoshift/header"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/probehandler"
	"github.com/acoshift/webstatic"
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

// New creates new app
func New() *hime.App {
	server := hime.New()
	server.ParseConfigFile("settings/routes.yaml")

	static := configfile.NewYAMLReader("static.yaml")
	server.TemplateFunc("static", func(s string) string {
		return "/-/" + static.StringDefault(s, s)
	})

	baseURL := config.String("base_url")
	server.Globals(hime.Globals{
		"baseURL": baseURL,
	})

	server.Template().
		Funcs(templateFunc()).
		ParseConfigFile("settings/template.yaml")

	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := httpmux.New()
	app.Mount(m)
	auth.Mount(m)
	editor.Mount(m)
	admin.Mount(m)

	h := middleware.Chain(
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

	mux := http.NewServeMux()
	// health check
	probe := probehandler.New()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ready") == "1" {
			// readiness probe
			probe.ServeHTTP(w, r)
			return
		}

		// liveness probe
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
		Dir:          "assets",
		CacheControl: "public, max-age=31536000, immutable",
	})))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "assets/favicon.ico") })

	mux.Handle("/", h)

	handler := middleware.Chain(
		errorLogger,
		defaultHeaders,
		middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{baseURL},
			IgnoreProto: true,
			ForbiddenHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Cross-site origin detected!", http.StatusForbidden)
			}),
		}),
	)(mux)

	server.Handler(handler)

	server.GracefulShutdown().
		Wait(5 * time.Second).
		Timeout(10 * time.Second).
		Notify(probe.Fail)

	return server
}

func turbolinks(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Turbolinks-Location", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func defaultHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.XContentTypeOptions, "nosniff")
		w.Header().Set(header.XXSSProtection, "1; mode=block")
		w.Header().Set(header.XFrameOptions, "deny")
		// w.Header().Set(header.ContentSecurityPolicy, "img-src https: data:; font-src https: data:; media-src https:;")
		w.Header().Set(header.CacheControl, "no-store")
		w.Header().Set(header.ReferrerPolicy, "same-origin")
		h.ServeHTTP(w, r)
	})
}

// errorLogger logs error and send error page back to response
func errorLogger(h http.Handler) http.Handler {
	errClient := config.ErrorClient()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				debug.PrintStack()

				if errClient != nil {
					nerr, _ := err.(error)
					errClient.Report(errorreporting.Entry{
						Error: nerr,
						Req:   r,
						Stack: debug.Stack(),
					})
				}
			}
		}()
		h.ServeHTTP(w, r)
	})
}
