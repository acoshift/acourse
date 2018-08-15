package app

import (
	"database/sql"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/acoshift/prefixhandler"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/goredis"
	"github.com/acoshift/webstatic"
	"github.com/go-redis/redis"
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
	redisClient  *redis.Client
	redisPrefix  string
	db           *sql.DB
)

// New creates new app
func New(config Config) http.Handler {
	auth = config.Auth
	loc = config.Location
	slackURL = config.SlackURL
	emailFrom = config.EmailFrom
	emailDialer = config.EmailDialer
	baseURL = config.BaseURL
	bucketHandle = config.BucketHandle
	bucketName = config.BucketName
	redisClient = config.RedisClient
	redisPrefix = config.RedisPrefix
	db = config.DB

	mux := http.NewServeMux()

	mux.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
		Dir:          "assets",
		CacheControl: "public, max-age=31536000",
	})))
	mux.Handle("/favicon.ico", fileHandler("assets/favicon.ico"))

	methodmux.FallbackHandler = hime.Handler(notFound)

	r := http.NewServeMux()

	r.Handle("/", methodmux.Get(
		hime.Handler(index),
	))

	// auth
	r.Handle("/signin", mustNotSignedIn(methodmux.GetPost(
		hime.Handler(signIn),
		hime.Handler(postSignIn),
	)))
	r.Handle("/signin/password", mustNotSignedIn(methodmux.GetPost(
		hime.Handler(signInPassword),
		hime.Handler(postSignInPassword),
	)))
	r.Handle("/signin/check-email", mustNotSignedIn(methodmux.Get(
		hime.Handler(checkEmail),
	)))
	r.Handle("/signin/link", mustNotSignedIn(methodmux.Get(
		hime.Handler(signInLink),
	)))
	r.Handle("/reset/password", mustNotSignedIn(methodmux.GetPost(
		hime.Handler(resetPassword),
		hime.Handler(postResetPassword),
	)))

	r.Handle("/openid", mustNotSignedIn(methodmux.Get(
		hime.Handler(openID),
	)))
	r.Handle("/openid/callback", mustNotSignedIn(methodmux.Get(
		hime.Handler(openIDCallback),
	)))
	r.Handle("/signup", mustNotSignedIn(methodmux.GetPost(
		hime.Handler(signUp),
		hime.Handler(postSignUp),
	)))
	r.Handle("/signout", methodmux.GetPost(
		hime.Handler(signOut),
		hime.Handler(signOut),
	)) // TODO: remove get signout

	// profile
	r.Handle("/profile", mustSignedIn(methodmux.Get(
		hime.Handler(profile),
	)))
	r.Handle("/profile/edit", mustSignedIn(methodmux.GetPost(
		hime.Handler(profileEdit),
		hime.Handler(postProfileEdit),
	)))

	// course
	{
		m := http.NewServeMux()
		m.Handle("/", methodmux.Get(
			hime.Handler(courseView),
		))
		m.Handle("/content", mustSignedIn(methodmux.Get(
			hime.Handler(courseContent),
		)))
		m.Handle("/enroll", mustSignedIn(methodmux.GetPost(
			hime.Handler(courseEnroll),
			hime.Handler(postCourseEnroll),
		)))
		m.Handle("/assignment", mustSignedIn(methodmux.Get(
			hime.Handler(courseAssignment),
		)))

		r.Handle("/course/", prefixhandler.New("/course", courseURLKey{}, m))
	}

	// editor
	{
		m := http.NewServeMux()
		m.Handle("/create", onlyInstructor(methodmux.GetPost(
			hime.Handler(editorCreate),
			hime.Handler(postEditorCreate),
		)))
		m.Handle("/course", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorCourse),
			hime.Handler(postEditorCourse),
		)))
		m.Handle("/content", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorContent),
			hime.Handler(postEditorContent),
		)))
		m.Handle("/content/create", isCourseOwner(methodmux.GetPost(
			hime.Handler(editorContentCreate),
			hime.Handler(postEditorContentCreate),
		)))
		// TODO: add middleware ?
		m.Handle("/content/edit", methodmux.GetPost(
			hime.Handler(editorContentEdit),
			hime.Handler(postEditorContentEdit),
		))

		r.Handle("/editor/", http.StripPrefix("/editor", m))
	}

	// admin
	{
		m := http.NewServeMux()
		m.Handle("/users", methodmux.Get(
			hime.Handler(adminUsers),
		))
		m.Handle("/courses", methodmux.Get(
			hime.Handler(adminCourses),
		))
		m.Handle("/payments/pending", methodmux.GetPost(
			hime.Handler(adminPendingPayments),
			hime.Handler(postAdminPendingPayment),
		))
		m.Handle("/payments/history", methodmux.Get(
			hime.Handler(adminHistoryPayments),
		))
		m.Handle("/payments/reject", methodmux.GetPost(
			hime.Handler(adminRejectPayment),
			hime.Handler(postAdminRejectPayment),
		))

		r.Handle("/admin/", onlyAdmin(http.StripPrefix("/admin", m)))
	}

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
				Client: config.RedisClient,
			}),
		}),
		fetchUser(),
	)(r))

	return mux
}

var notFoundImages = []string{
	"https://storage.googleapis.com/acourse/static/9961f3c1-575f-4b98-af4f-447566ee1cb3.png",
	"https://storage.googleapis.com/acourse/static/b14a40c9-d3a4-465d-9453-ce7fcfbc594c.png",
}

func notFound(ctx *hime.Context) error {
	page := newPage(ctx)
	page["Image"] = notFoundImages[rand.Intn(len(notFoundImages))]
	ctx.ResponseWriter().Header().Set(header.XContentTypeOptions, "nosniff")
	return ctx.Status(http.StatusNotFound).View("error.not-found", page)
}
