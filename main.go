package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/profiler"
	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/acoshift/middleware"
	"github.com/acoshift/probehandler"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/app"
	"github.com/acoshift/acourse/internal"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// init profiler
	profiler.Start(profiler.Config{Service: "acourse"})

	firApp, err := firebase.InitializeApp(ctx, firebase.AppOptions{
		ProjectID: config.String("project_id"),
	}, option.WithCredentialsFile("config/service_account"))
	if err != nil {
		log.Fatal(err)
	}
	firAuth := firApp.Auth()

	// init google storage
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("config/service_account"))
	if err != nil {
		log.Fatal(err)
	}
	bucketHandle := storageClient.Bucket(config.String("bucket"))

	// init email client
	emailDialer := gomail.NewPlainDialer(config.String("email_server"), config.Int("email_port"), config.String("email_user"), config.String("email_password"))

	// init redis pool
	redisClient := redis.NewClient(&redis.Options{
		MaxRetries:  3,
		PoolSize:    5,
		IdleTimeout: 60 * time.Minute,
		Addr:        config.String("redis_addr"),
		Password:    config.String("redis_pass"),
	})

	// init databases
	db, err := sql.Open("postgres", config.String("sql_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)

	himeApp := hime.New()
	himeApp.ParseConfigFile("settings/server.yaml")
	himeApp.ParseConfigFile("settings/routes.yaml")

	static := configfile.NewYAMLReader("static.yaml")
	himeApp.TemplateFunc("static", func(s string) string {
		return "/-/" + static.StringDefault(s, s)
	})

	baseURL := config.String("base_url")

	himeApp.Template().
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

	mux.Handle("/", app.New(app.Config{
		DB:            db,
		BaseURL:       baseURL,
		RedisClient:   redisClient,
		RedisPrefix:   config.String("redis_prefix"),
		SessionSecret: config.Bytes("session_secret"),
		Auth:          firAuth,
		Location:      loc,
		SlackURL:      config.String("slack_url"),
		EmailFrom:     config.String("email_from"),
		EmailDialer:   emailDialer,
		BucketHandle:  bucketHandle,
		BucketName:    config.String("bucket"),
	}))

	h := middleware.Chain(
		internal.ErrorRecovery,
		internal.SetHeaders,
		middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{baseURL},
			IgnoreProto: true,
			ForbiddenHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Cross-site origin detected!", http.StatusForbidden)
			}),
		}),
	)(mux)

	err = himeApp.
		Handler(h).
		GracefulShutdown().
		Notify(probe.Fail).
		ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
