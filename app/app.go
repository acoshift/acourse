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
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/prefixhandler"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/acoshift/webstatic"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
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
			"openid.google":          "/openid?p=google.com",
			"openid.github":          "/openid?p=github.com",
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

		methodmux.FallbackHandler = hime.H(notFound)

		r := http.NewServeMux()

		r.Handle("/", methodmux.Get(
			hime.H(index),
		))

		// auth
		r.Handle("/signin", mustNotSignedIn(methodmux.GetPost(
			hime.H(signIn),
			hime.H(postSignIn),
		)))
		r.Handle("/signin/password", mustNotSignedIn(methodmux.GetPost(
			hime.H(signInPassword),
			hime.H(postSignInPassword),
		)))
		r.Handle("/signin/check-email", mustNotSignedIn(methodmux.Get(
			hime.H(checkEmail),
		)))
		r.Handle("/signin/link", mustNotSignedIn(methodmux.Get(
			hime.H(signInLink),
		)))
		r.Handle("/reset/password", mustNotSignedIn(methodmux.GetPost(
			hime.H(resetPassword),
			hime.H(postResetPassword),
		)))

		r.Handle("/openid", mustNotSignedIn(methodmux.Get(
			hime.H(openID),
		)))
		r.Handle("/openid/callback", mustNotSignedIn(methodmux.Get(
			hime.H(openIDCallback),
		)))
		r.Handle("/signup", mustNotSignedIn(methodmux.GetPost(
			hime.H(signUp),
			hime.H(postSignUp),
		)))
		r.Handle("/signout", methodmux.GetPost(
			hime.H(signOut),
			hime.H(signOut),
		)) // TODO: remove get signout

		// profile
		r.Handle("/profile", mustSignedIn(methodmux.Get(
			hime.H(profile),
		)))
		r.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
			hime.H(profileEdit),
			hime.H(postProfileEdit),
		)))

		// course
		courseMux := http.NewServeMux()
		courseMux.Handle("/", methodmux.Get(
			hime.H(courseView),
		))
		courseMux.Handle("/content", mustSignedIn(methodmux.Get(
			hime.H(courseContent),
		)))
		courseMux.Handle("/enroll", mustSignedIn(methodmux.GetPost(
			hime.H(courseEnroll),
			hime.H(postCourseEnroll),
		)))
		courseMux.Handle("/assignment", mustSignedIn(methodmux.Get(
			hime.H(courseAssignment),
		)))

		r.Handle("/course/", prefixhandler.New("/course", courseURLKey{}, courseMux))

		// editor
		r.Handle("/editor/create", onlyInstructor(methodmux.GetPost(
			hime.H(editorCreate),
			hime.H(postEditorCreate),
		)))
		r.Handle("/editor/course", isCourseOwner(methodmux.GetPost(
			hime.H(editorCourse),
			hime.H(postEditorCourse),
		)))
		r.Handle("/editor/content", isCourseOwner(methodmux.GetPost(
			hime.H(editorContent),
			hime.H(postEditorContent),
		)))
		r.Handle("/editor/content/create", isCourseOwner(methodmux.GetPost(
			hime.H(editorContentCreate),
			hime.H(postEditorContentCreate),
		)))
		r.Handle("/editor/content/edit", methodmux.GetPost(
			hime.H(editorContentEdit),
			hime.H(postEditorContentEdit),
		))

		// admin
		r.Handle("/admin/users", onlyAdmin(methodmux.Get(
			hime.H(adminUsers),
		)))
		r.Handle("/admin/courses", onlyAdmin(methodmux.Get(
			hime.H(adminCourses),
		)))
		r.Handle("/admin/payments/pending", onlyAdmin(methodmux.GetPost(
			hime.H(adminPendingPayments),
			hime.H(postAdminPendingPayment),
		)))
		r.Handle("/admin/payments/history", onlyAdmin(methodmux.Get(
			hime.H(adminHistoryPayments),
		)))
		r.Handle("/admin/payments/reject", onlyAdmin(methodmux.GetPost(
			hime.H(adminRejectPayment),
			hime.H(postAdminRejectPayment),
		)))

		mux.Handle("/", middleware.Chain(
			session.Middleware(session.Config{
				Secret:   config.SessionSecret,
				Path:     "/",
				MaxAge:   7 * 24 * time.Hour,
				HTTPOnly: true,
				Secure:   session.PreferSecure,
				SameSite: session.SameSiteLax,
				Rolling:  true,
				Proxy:    true,
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
