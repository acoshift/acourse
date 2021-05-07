package app

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"cloud.google.com/go/errorreporting"
	"github.com/acoshift/configfile"
	"github.com/acoshift/header"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
	"github.com/acoshift/probehandler"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/httpmux"
	"github.com/moonrhythm/session"
	"github.com/moonrhythm/session/store"
	"github.com/moonrhythm/webstatic/v4"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/acoshift/acourse/internal/app/admin"
	"github.com/acoshift/acourse/internal/app/app"
	"github.com/acoshift/acourse/internal/app/auth"
	"github.com/acoshift/acourse/internal/app/editor"
	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
)

// Config is the App's config
type Config struct {
	Routes         []byte
	Static         []byte
	Template       fs.FS
	TemplateConfig []byte
	Assets         fs.FS
}

// New creates new app
func New(cfg Config) *hime.App {
	server := hime.New()
	server.ETag = true
	server.ParseConfig(cfg.Routes)

	static := configfile.NewYAMLReaderFromReader(bytes.NewReader(cfg.Static))
	server.TemplateFunc("static", func(s string) string {
		return "/-/" + static.StringDefault(s, s)
	})

	baseURL := config.String("base_url")
	server.Globals(hime.Globals{
		"baseURL": baseURL,
	})

	server.Template().
		FS(cfg.Template).
		Funcs(templateFunc()).
		ParseConfig(cfg.TemplateConfig)

	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := httpmux.New()
	app.Mount(m)
	auth.Mount(m)
	editor.Mount(m)
	admin.Mount(m)

	h := middleware.Chain(
		traceMiddleware,
		pgctx.Middleware(config.DBClient()),
		redisctx.Middleware(config.RedisClient(), config.String("redis_prefix")),
		session.Middleware(session.Config{
			Secret:      config.Bytes("session_secret"),
			Path:        "/",
			MaxAge:      7 * 24 * time.Hour,
			HTTPOnly:    true,
			Secure:      session.PreferSecure,
			SameSite:    http.SameSiteLaxMode,
			Rolling:     true,
			Proxy:       true,
			Resave:      true,
			ResaveAfter: 24 * time.Hour,
			Store: (&store.SQL{
				DB: config.DBClient(),
			}).GeneratePostgreSQLStatement("sessions", false).GCEvery(time.Hour),
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

	mux.Handle("/-/", http.StripPrefix("/-", &webstatic.Handler{
		FileSystem:   http.FS(cfg.Assets),
		CacheControl: "public, max-age=31536000, immutable",
	}))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "assets/favicon.ico") })

	mux.Handle("/", h)

	handler := middleware.Chain(
		errorLogger,
		defaultHeaders,
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
		w.Header().Set(header.CacheControl, "no-cache")
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
				if e, ok := err.(error); ok {
					if errors.Is(e, context.Canceled) {
						return
					}
					if errors.Is(e, driver.ErrBadConn) {
						return
					}
					if pgsql.IsQueryCanceled(e) {
						return
					}
				}

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

func traceMiddleware(h http.Handler) http.Handler {
	return &ochttp.Handler{
		Handler: h,
		FormatSpanName: func(r *http.Request) string {
			proto := r.Header.Get("X-Forwarded-Proto")
			return proto + "://" + r.Host + r.RequestURI
		},
		StartOptions: trace.StartOptions{
			Sampler:  trace.AlwaysSample(),
			SpanKind: trace.SpanKindServer,
		},
	}
}
