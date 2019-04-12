package main

import (
	"database/sql"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"github.com/acoshift/configfile"
	"github.com/acoshift/probehandler"
	"github.com/acoshift/webstatic"
	_ "github.com/lib/pq"
	"github.com/moonrhythm/hime"
	redisstore "github.com/moonrhythm/session/store/goredis"

	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/middleware"
)

func main() {
	// init databases
	db, err := sql.Open("postgres", config.String("sql_url"))
	must(err)
	defer db.Close()
	db.SetMaxOpenConns(config.IntDefault("sql_max_open_conns", 5))

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
		Funcs(app.TemplateFunc()).
		ParseConfigFile("settings/template.yaml")

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

	mux.Handle("/", app.Handler(app.Config{
		DB:            db,
		SessionSecret: config.Bytes("session_secret"),
		SessionStore: redisstore.New(redisstore.Config{
			Prefix: config.String("redis_prefix"),
			Client: config.RedisClient(),
		}),
		RedisClient: config.RedisClient(),
		RedisPrefix: config.String("redis_prefix"),
	}))

	h := middleware.Chain(
		middleware.ErrorLogger(config.ErrorClient()),
		middleware.SetHeaders,
		middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{baseURL},
			IgnoreProto: true,
			ForbiddenHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Cross-site origin detected!", http.StatusForbidden)
			}),
		}),
	)(mux)

	server.GracefulShutdown().
		Wait(5 * time.Second).
		Timeout(10 * time.Second).
		Notify(probe.Fail)

	err = server.
		Handler(h).
		Address(":8080").
		ListenAndServe()
	must(err)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
