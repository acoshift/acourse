package app

import (
	"database/sql"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/csrf"
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

		r.Get("/", hime.H(index))

		// auth
		r.Get(app.Route("signin"), mustNotSignedIn(hime.H(signIn)))
		r.Post(app.Route("signin"), mustNotSignedIn(hime.H(postSignIn)))
		r.Get(app.Route("signin.password"), mustNotSignedIn(hime.H(signInPassword)))
		r.Post(app.Route("signin.password"), mustNotSignedIn(hime.H(postSignInPassword)))
		r.Get(app.Route("signin.check-email"), mustNotSignedIn(hime.H(checkEmail)))
		r.Get(app.Route("signin.link"), mustNotSignedIn(hime.H(signInLink)))
		r.Get(app.Route("reset.password"), mustNotSignedIn(hime.H(resetPassword)))
		r.Post(app.Route("reset.password"), mustNotSignedIn(hime.H(postResetPassword)))
		r.Get(app.Route("openid"), mustNotSignedIn(hime.H(openID)))
		r.Get(app.Route("openid.callback"), mustNotSignedIn(hime.H(openIDCallback)))
		r.Get(app.Route("signup"), mustNotSignedIn(hime.H(signUp)))
		r.Post(app.Route("signup"), mustNotSignedIn(hime.H(postSignUp)))
		r.Get(app.Route("signout"), hime.H(signOut)) // TODO: remove get signout
		r.Post(app.Route("signout"), hime.H(signOut))

		// profile
		r.Get(app.Route("profile"), mustSignedIn(hime.H(profile)))
		r.Get(app.Route("profile.edit"), mustSignedIn(hime.H(profileEdit)))
		r.Post(app.Route("profile.edit"), mustSignedIn(hime.H(postProfileEdit)))

		// course
		r.Get(app.Route("course", ":courseURL"), hime.H(courseView))
		r.Get(app.Route("course", ":courseURL", "content"), hime.H(courseContent))
		r.Get(app.Route("course", ":courseURL", "enroll"), hime.H(courseEnroll))
		r.Post(app.Route("course", ":courseURL", "enroll"), hime.H(postCourseEnroll))
		r.Get(app.Route("course", ":courseURL", "assignment"), hime.H(courseAssignment))

		// editor
		r.Get(app.Route("editor.create"), onlyInstructor(hime.H(editorCreate)))
		r.Post(app.Route("editor.create"), onlyInstructor(hime.H(postEditorCreate)))
		r.Get(app.Route("editor.course"), isCourseOwner(hime.H(editorCourse)))
		r.Post(app.Route("editor.course"), isCourseOwner(hime.H(postEditorCourse)))
		r.Get(app.Route("editor.content"), isCourseOwner(hime.H(editorContent)))
		r.Post(app.Route("editor.content"), isCourseOwner(hime.H(postEditorContent)))
		r.Get(app.Route("editor.content.create"), isCourseOwner(hime.H(editorContentCreate)))
		r.Post(app.Route("editor.content.create"), isCourseOwner(hime.H(postEditorContentCreate)))
		r.Get(app.Route("editor.content.edit"), hime.H(editorContentEdit))
		r.Post(app.Route("editor.content.edit"), hime.H(postEditorContentEdit))

		// admin
		r.Get(app.Route("admin.users"), onlyAdmin(hime.H(adminUsers)))
		r.Get(app.Route("admin.courses"), onlyAdmin(hime.H(adminCourses)))
		r.Get(app.Route("admin.payments.pending"), onlyAdmin(hime.H(adminPendingPayments)))
		r.Post(app.Route("admin.payments.pending"), onlyAdmin(hime.H(postAdminPendingPayment)))
		r.Get(app.Route("admin.payments.history"), onlyAdmin(hime.H(adminHistoryPayments)))
		r.Get(app.Route("admin.payments.reject"), onlyAdmin(hime.H(adminRejectPayment)))
		r.Post(app.Route("admin.payments.reject"), onlyAdmin(hime.H(postAdminRejectPayment)))

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
			csrf.New(csrf.Config{
				Origins: []string{config.BaseURL},
				ForbiddenHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Cross-site origin detected!", http.StatusForbidden)
				}),
			}),
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
