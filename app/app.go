package app

import (
	"database/sql"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/acourse/view"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/acoshift/webstatic"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/gomail.v2"
)

var (
	auth         *firebase.Auth
	loc          *time.Location
	slackURL     string
	emailFrom    string
	emailDialer  *gomail.Dialer
	baseURL      string
	bucketHandle *storage.BucketHandle
	bucketName   string
	redisPool    *redis.Pool
	redisPrefix  string
	cachePool    *redis.Pool
	cachePrefix  string
	db           *sql.DB
)

// New creates new app
func New(config Config) hime.HandlerFactory {
	auth = config.Auth
	loc = config.Location
	slackURL = config.SlackURL
	emailFrom = config.EmailFrom
	emailDialer = config.EmailDialer
	baseURL = config.BaseURL
	bucketHandle = config.BucketHandle
	bucketName = config.BucketName
	redisPool = config.RedisPool
	redisPrefix = config.RedisPrefix
	cachePool = config.CachePool
	cachePrefix = config.CachePrefix
	db = config.DB

	return func(app hime.App) http.Handler {
		mux := http.NewServeMux()

		app.Routes(hime.Routes{
			"index":              "/",
			"signin":             "/signin",
			"signin.password":    "/signin/password",
			"signin.check-email": "/signin/check-email",
			"signin.link":        "/signin/link",
			"openid":             "/openid",
			"openid.callback":    "/openid/callback",
			"signup":             "/signup",
			"signout":            "/signout",
			"reset.password":     "/reset/password",
		})

		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		editor := httprouter.New()
		editor.HandleMethodNotAllowed = false
		editor.HandleOPTIONS = false
		editor.NotFound = http.HandlerFunc(notFound)
		editor.Get("/create", onlyInstructor(http.HandlerFunc(editorCreate)))
		editor.Post("/create", onlyInstructor(http.HandlerFunc(postEditorCreate)))
		editor.Get("/course", isCourseOwner(http.HandlerFunc(editorCourse)))
		editor.Post("/course", isCourseOwner(http.HandlerFunc(postEditorCourse)))
		editor.Get("/content", isCourseOwner(http.HandlerFunc(editorContent)))
		editor.Post("/content", isCourseOwner(http.HandlerFunc(postEditorContent)))
		editor.Get("/content/create", isCourseOwner(http.HandlerFunc(editorContentCreate)))
		editor.Post("/content/create", isCourseOwner(http.HandlerFunc(postEditorContentCreate)))
		editor.Get("/content/edit", http.HandlerFunc(editorContentEdit))
		editor.Post("/content/edit", http.HandlerFunc(postEditorContentEdit))

		admin := httprouter.New()
		admin.HandleMethodNotAllowed = false
		admin.HandleOPTIONS = false
		admin.NotFound = http.HandlerFunc(notFound)
		admin.Get("/users", http.HandlerFunc(adminUsers))
		admin.Get("/courses", http.HandlerFunc(adminCourses))
		admin.Get("/payments/pending", http.HandlerFunc(adminPendingPayments))
		admin.Post("/payments/pending", http.HandlerFunc(postAdminPendingPayment))
		admin.Get("/payments/history", http.HandlerFunc(adminHistoryPayments))
		admin.Get("/payments/reject", http.HandlerFunc(adminRejectPayment))
		admin.Post("/payments/reject", http.HandlerFunc(postAdminRejectPayment))

		main := http.NewServeMux()
		main.Handle("/", http.HandlerFunc(index))
		main.Handle(app.Route("signin"), mustNotSignedIn(http.HandlerFunc(signIn)))
		main.Handle(app.Route("signin.password"), mustNotSignedIn(http.HandlerFunc(signInPassword)))
		main.Handle(app.Route("signin.check-email"), mustNotSignedIn(http.HandlerFunc(checkEmail)))
		main.Handle(app.Route("signin.link"), mustNotSignedIn(http.HandlerFunc(signInLink)))
		main.Handle(app.Route("openid"), mustNotSignedIn(http.HandlerFunc(openID)))
		main.Handle(app.Route("openid.callback"), mustNotSignedIn(http.HandlerFunc(openIDCallback)))
		main.Handle(app.Route("signup"), mustNotSignedIn(http.HandlerFunc(signUp)))
		main.Handle(app.Route("signout"), http.HandlerFunc(signOut))
		main.Handle(app.Route("reset.password"), mustNotSignedIn(http.HandlerFunc(resetPassword)))
		main.Handle("/profile", mustSignedIn(http.HandlerFunc(profile)))
		main.Handle("/profile/edit", mustSignedIn(http.HandlerFunc(profileEdit)))
		main.Handle("/course/", http.StripPrefix("/course/", courseHandler()))
		main.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
		main.Handle("/editor/", http.StripPrefix("/editor", editor))

		mux.Handle("/~/", http.StripPrefix("/~", cache(webstatic.New("static"))))
		mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

		mux.Handle("/", middleware.Chain(
			session.Middleware(session.Config{
				Secret:   config.SessionSecret,
				Path:     "/",
				MaxAge:   30 * 24 * time.Hour,
				HTTPOnly: true,
				Secure:   session.PreferSecure,
				SameSite: session.SameSiteLax,
				Store: redisstore.New(redisstore.Config{
					Prefix: config.RedisPrefix,
					Pool:   config.RedisPool,
				}),
			}),
			fetchUser(),
			csrf(config.BaseURL, config.XSRFSecret),
		)(main))

		return middleware.Chain(
			errorRecovery,
			setHeaders,
		)(mux)
	}
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	view.NotFound(w, r)
}
