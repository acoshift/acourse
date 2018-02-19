package app

import (
	"database/sql"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/acoshift/webstatic"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/gomail.v2"
	yaml "gopkg.in/yaml.v2"
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
	staticConf   = make(map[string]string)
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

	// load static config
	// TODO: move to main
	{
		bs, _ := ioutil.ReadFile("static.yaml")
		yaml.Unmarshal(bs, &staticConf)
	}

	return func(app hime.App) http.Handler {
		loadTemplates(app)

		app.Routes(hime.Routes{
			"index":                  "/",
			"signin":                 "/signin",
			"signin.password":        "/signin/password",
			"signin.check-email":     "/signin/check-email",
			"signin.link":            "/signin/link",
			"openid":                 "/openid",
			"openid.callback":        "/openid/callback",
			"signup":                 "/signup",
			"signout":                "/signout",
			"reset.password":         "/reset/password",
			"profile":                "/profile",
			"profile.edit":           "/profile/edit",
			"course":                 "/course/",
			"editor.create":          "/editor/create",
			"editor.course":          "/editor/course",
			"editor.content":         "/editor/content",
			"editor.content.create":  "/editor/content/create",
			"editor.content.edit":    "/editor/content/edit",
			"admin.users":            "/admin/users",
			"admin.courses":          "/admin/courses",
			"admin.payments.pending": "/admin/payments/pending",
			"admin.payments.history": "/admin/payments/history",
			"admin.payments.reject":  "/admin/payments/reject",
		})

		mux := http.NewServeMux()

		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mux.Handle("/~/", http.StripPrefix("/~", cache(webstatic.New("static"))))
		mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

		r := httprouter.New()
		r.HandleMethodNotAllowed = false
		r.HandleOPTIONS = false
		r.NotFound = hime.H(notFound)

		r.Get("/", http.HandlerFunc(index))

		// auth
		r.Get(app.Route("signin"), mustNotSignedIn(http.HandlerFunc(signIn)))
		r.Post(app.Route("signin"), mustNotSignedIn(http.HandlerFunc(postSignIn)))
		r.Get(app.Route("signin.password"), mustNotSignedIn(http.HandlerFunc(signInPassword)))
		r.Post(app.Route("signin.password"), mustNotSignedIn(http.HandlerFunc(postSignInPassword)))
		r.Get(app.Route("signin.check-email"), mustNotSignedIn(http.HandlerFunc(checkEmail)))
		r.Get(app.Route("signin.link"), mustNotSignedIn(http.HandlerFunc(signInLink)))
		r.Get(app.Route("reset.password"), mustNotSignedIn(http.HandlerFunc(resetPassword)))
		r.Post(app.Route("reset.password"), mustNotSignedIn(http.HandlerFunc(postResetPassword)))
		r.Get(app.Route("openid"), mustNotSignedIn(http.HandlerFunc(openID)))
		r.Get(app.Route("openid.callback"), mustNotSignedIn(http.HandlerFunc(openIDCallback)))
		r.Get(app.Route("signup"), mustNotSignedIn(http.HandlerFunc(signUp)))
		r.Post(app.Route("signup"), mustNotSignedIn(http.HandlerFunc(postSignUp)))
		r.Get(app.Route("signout"), http.HandlerFunc(signOut)) // TODO: remove get signout
		r.Post(app.Route("signout"), http.HandlerFunc(signOut))

		// profile
		r.Get(app.Route("profile"), mustSignedIn(http.HandlerFunc(profile)))
		r.Get(app.Route("profile.edit"), mustSignedIn(http.HandlerFunc(profileEdit)))
		r.Post(app.Route("profile.edit"), mustSignedIn(http.HandlerFunc(postProfileEdit)))

		// course
		r.Get(app.Route("course", ":courseURL"), http.HandlerFunc(courseView))
		r.Get(app.Route("course", ":courseURL", "content"), http.HandlerFunc(courseContent))
		r.Get(app.Route("course", ":courseURL", "enroll"), http.HandlerFunc(courseEnroll))
		r.Post(app.Route("course", ":courseURL", "enroll"), http.HandlerFunc(postCourseEnroll))
		r.Get(app.Route("course", ":courseURL", "assignment"), http.HandlerFunc(courseAssignment))

		// editor
		r.Get(app.Route("editor.create"), onlyInstructor(http.HandlerFunc(editorCreate)))
		r.Post(app.Route("editor.create"), onlyInstructor(http.HandlerFunc(postEditorCreate)))
		r.Get(app.Route("editor.course"), isCourseOwner(http.HandlerFunc(editorCourse)))
		r.Post(app.Route("editor.course"), isCourseOwner(http.HandlerFunc(postEditorCourse)))
		r.Get(app.Route("editor.content"), isCourseOwner(http.HandlerFunc(editorContent)))
		r.Post(app.Route("editor.content"), isCourseOwner(http.HandlerFunc(postEditorContent)))
		r.Get(app.Route("editor.content.create"), isCourseOwner(http.HandlerFunc(editorContentCreate)))
		r.Post(app.Route("editor.content.create"), isCourseOwner(http.HandlerFunc(postEditorContentCreate)))
		r.Get(app.Route("editor.content.edit"), http.HandlerFunc(editorContentEdit))
		r.Post(app.Route("editor.content.edit"), http.HandlerFunc(postEditorContentEdit))

		// admin
		r.Get(app.Route("admin.users"), onlyAdmin(http.HandlerFunc(adminUsers)))
		r.Get(app.Route("admin.courses"), onlyAdmin(http.HandlerFunc(adminCourses)))
		r.Get(app.Route("admin.payments.pending"), onlyAdmin(http.HandlerFunc(adminPendingPayments)))
		r.Post(app.Route("admin.payments.pending"), onlyAdmin(http.HandlerFunc(postAdminPendingPayment)))
		r.Get(app.Route("admin.payments.history"), onlyAdmin(http.HandlerFunc(adminHistoryPayments)))
		r.Get(app.Route("admin.payments.reject"), onlyAdmin(http.HandlerFunc(adminRejectPayment)))
		r.Post(app.Route("admin.payments.reject"), onlyAdmin(http.HandlerFunc(postAdminRejectPayment)))

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
		)(r))

		return middleware.Chain(
			errorRecovery,
			setHeaders,
		)(mux)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

func notFound(ctx hime.Context) hime.Result {
	page := newPage(ctx)
	page["Image"] = notFoundImages[rand.Intn(len(notFoundImages))]
	ctx.ResponseWriter().Header().Set(header.XContentTypeOptions, "nosniff")
	return ctx.Status(http.StatusNotFound).View("error.not-found", page)
}
