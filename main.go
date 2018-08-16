package main

import (
	"context"
	"database/sql"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/profiler"
	"cloud.google.com/go/storage"
	"github.com/acoshift/configfile"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/probehandler"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/goredis"
	"github.com/acoshift/webstatic"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/redisctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/admin"
	"github.com/acoshift/acourse/controller/app"
	"github.com/acoshift/acourse/controller/auth"
	"github.com/acoshift/acourse/controller/editor"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/email"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/internal"
	"github.com/acoshift/acourse/notify"
	"github.com/acoshift/acourse/service"
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

	// init email sender
	emailSender := email.NewSMTPSender(email.SMTPConfig{
		Server:   config.String("email_server"),
		Port:     config.Int("email_port"),
		User:     config.String("email_user"),
		Password: config.String("email_password"),
		From:     config.String("email_from"),
	})

	adminNotifier := notify.NewOutgoingWebhookAdminNotifier(config.String("slack_url"))

	// init redis pool
	redisClient := redis.NewClient(&redis.Options{
		MaxRetries:  3,
		PoolSize:    5,
		IdleTimeout: 60 * time.Minute,
		Addr:        config.String("redis_addr"),
		Password:    config.String("redis_pass"),
	})

	fileStorage := file.NewGCS(storageClient, config.String("bucket"))
	imageResizeEncoder := image.NewJPEGResizeEncoder()

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
		Funcs(internal.TemplateFunc(loc)).
		ParseConfigFile("settings/template.yaml")

	methodmux.FallbackHandler = hime.Handler(share.NotFound)

	svc := service.New(service.Config{
		Auth:               firAuth,
		EmailSender:        emailSender,
		BaseURL:            baseURL,
		FileStorage:        fileStorage,
		ImageResizeEncoder: imageResizeEncoder,
		MagicLinkCallback:  himeApp.Route("auth.signin.link"),
		OpenIDCallback:     himeApp.Route("auth.openid.callback"),
	})

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
		CacheControl: "public, max-age=31536000",
	})))

	mux.Handle("/favicon.ico", internal.FileHandler("assets/favicon.ico"))

	m := http.NewServeMux()

	m.Handle("/", app.New(app.Config{
		BaseURL:            baseURL,
		Auth:               firAuth,
		AdminNotifier:      adminNotifier,
		FileStorage:        fileStorage,
		ImageResizeEncoder: imageResizeEncoder,
	}))

	m.Handle("/auth/", http.StripPrefix("/auth", middleware.Chain(
		internal.NotSignedIn,
	)(auth.New(auth.Config{
		Service: svc,
	}))))

	m.Handle("/editor/", http.StripPrefix("/editor", middleware.Chain(
	// -
	)(editor.New(editor.Config{
		FileStorage:        fileStorage,
		ImageResizeEncoder: imageResizeEncoder,
	}))))

	m.Handle("/admin/", http.StripPrefix("/admin", middleware.Chain(
		internal.OnlyAdmin,
	)(admin.New(admin.Config{
		Location:    loc,
		EmailSender: emailSender,
	}))))

	mux.Handle("/", middleware.Chain(
		sqlctx.Middleware(db),
		redisctx.Middleware(redisClient, config.String("redis_prefix")),
		session.Middleware(session.Config{
			Secret:   config.Bytes("session_secret"),
			Path:     "/",
			MaxAge:   7 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			SameSite: session.SameSiteLax,
			Rolling:  true,
			Proxy:    true,
			Store: redisstore.New(redisstore.Config{
				Prefix: config.String("redis_prefix"),
				Client: redisClient,
			}),
		}),
		appctx.Middleware,
	)(m))

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
