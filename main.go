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
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/moonrhythm/hime"
	redisstore "github.com/moonrhythm/session/store/goredis"

	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/middleware"
	"github.com/acoshift/acourse/internal/service/admin"
)

func main() {
	time.Local = time.UTC

	loc, err := time.LoadLocation(config.StringDefault("location", "Asia/Bangkok"))
	must(err)

	// init redis pool
	redisClient := redis.NewClient(&redis.Options{
		MaxRetries:  config.IntDefault("redis_max_retries", 3),
		PoolSize:    config.IntDefault("redis_pool_size", 5),
		IdleTimeout: config.DurationDefault("redis_idle_timeout", 60*time.Minute),
		Addr:        config.String("redis_addr"),
		Password:    config.String("redis_pass"),
	})

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
		"baseURL":  baseURL,
		"location": loc,
	})

	server.Template().
		Funcs(app.TemplateFunc(loc)).
		ParseConfigFile("settings/template.yaml")

	admin.Init()

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
			Client: redisClient,
		}),
		RedisClient: redisClient,
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
