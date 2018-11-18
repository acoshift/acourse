package main

import (
	"context"
	"database/sql"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	firadmin "github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/middleware"
	"github.com/acoshift/probehandler"
	"github.com/acoshift/webstatic"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/moonrhythm/hime"
	"github.com/moonrhythm/session"
	redisstore "github.com/moonrhythm/session/store/goredis"
	"google.golang.org/api/option"

	"github.com/acoshift/acourse/internal"
	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/context/redisctx"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/service/admin"
	"github.com/acoshift/acourse/internal/service/auth"
	"github.com/acoshift/acourse/internal/service/course"
	"github.com/acoshift/acourse/internal/service/email"
	"github.com/acoshift/acourse/internal/service/file"
	"github.com/acoshift/acourse/internal/service/firebase"
	"github.com/acoshift/acourse/internal/service/image"
	"github.com/acoshift/acourse/internal/service/notify"
	"github.com/acoshift/acourse/internal/service/payment"
	"github.com/acoshift/acourse/internal/service/user"
)

func main() {
	time.Local = time.UTC

	config := configfile.NewReader("config")

	loc, err := time.LoadLocation(config.StringDefault("location", "Asia/Bangkok"))
	must(err)

	ctx := context.Background()

	googClientOpts := []option.ClientOption{option.WithCredentialsFile("config/service_account")}

	serviceName := config.StringDefault("service", "acourse")
	projectID := config.String("project_id")

	// init error reporting, ignore error
	errClient, _ := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Println(err)
		},
	}, googClientOpts...)

	firApp, err := firadmin.InitializeApp(ctx, firadmin.AppOptions{
		ProjectID: projectID,
	}, googClientOpts...)
	must(err)
	firAuth := firApp.Auth()

	// init google storage
	storageClient, err := storage.NewClient(ctx, googClientOpts...)
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

	himeApp := hime.New()
	himeApp.ParseConfigFile("settings/server.yaml")
	himeApp.ParseConfigFile("settings/routes.yaml")

	static := configfile.NewYAMLReader("static.yaml")
	himeApp.TemplateFunc("static", func(s string) string {
		return "/-/" + static.StringDefault(s, s)
	})

	baseURL := config.String("base_url")

	himeApp.Template().
		Funcs(internal.TemplateFunc(loc)).
		ParseConfigFile("settings/template.yaml")

	// init services
	email.InitSMTP(email.SMTPConfig{
		Server:   config.String("email_server"),
		Port:     config.Int("email_port"),
		User:     config.String("email_user"),
		Password: config.String("email_password"),
		From:     config.String("email_from"),
	})
	file.InitGCS(storageClient, config.String("bucket"))
	image.Init()
	firebase.Init(firAuth)
	notify.Init(config.String("slack_url"))
	auth.Init(baseURL, himeApp.Route("auth.openid.callback"))
	user.Init()
	course.Init()
	payment.Init()
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

	mux.Handle("/favicon.ico", internal.FileHandler("assets/favicon.ico"))

	mux.Handle("/", middleware.Chain(
		sqlctx.Middleware(db),
		redisctx.Middleware(redisClient, config.String("redis_prefix")),
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
				Client: redisClient,
			}),
		}),
		internal.Turbolinks,
		appctx.Middleware(),
	)(app.Handler(baseURL, loc)))

	h := middleware.Chain(
		internal.ErrorLogger(errClient),
		internal.SetHeaders,
		middleware.CSRF(middleware.CSRFConfig{
			Origins:     []string{baseURL},
			IgnoreProto: true,
			ForbiddenHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Cross-site origin detected!", http.StatusForbidden)
			}),
		}),
	)(mux)

	himeApp.GracefulShutdown().
		Notify(probe.Fail)

	err = himeApp.
		Handler(h).
		ListenAndServe()
	must(err)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
