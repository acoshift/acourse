package app

import (
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
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
	redisPool = config.RedisPool
	redisPrefix = config.RedisPrefix
	cachePool = config.CachePool
	cachePrefix = config.CachePrefix

	// create middlewares
	isCourseOwner := isCourseOwner(config.DB)

	// create mux
	mux := http.NewServeMux()

	editor := http.NewServeMux()
	editor.Handle("/create", onlyInstructor(http.HandlerFunc(editorCreate)))
	editor.Handle("/course", isCourseOwner(http.HandlerFunc(editorCourse)))
	editor.Handle("/content", isCourseOwner(http.HandlerFunc(editorContent)))
	editor.Handle("/content/create", isCourseOwner(http.HandlerFunc(editorContentCreate)))
	editor.Handle("/content/edit", http.HandlerFunc(editorContentEdit))

	admin := http.NewServeMux()
	admin.Handle("/users", http.HandlerFunc(adminUsers))
	admin.Handle("/courses", http.HandlerFunc(adminCourses))
	admin.Handle("/payments/pending", http.HandlerFunc(adminPendingPayments))
	admin.Handle("/payments/history", http.HandlerFunc(adminHistoryPayments))
	admin.Handle("/payments/reject", http.HandlerFunc(adminRejectPayment))

	main := http.NewServeMux()
	main.Handle("/", http.HandlerFunc(index))
	main.Handle("/signin", mustNotSignedIn(http.HandlerFunc(signIn)))
	main.Handle("/signin/password", mustNotSignedIn(http.HandlerFunc(signInPassword)))
	main.Handle("/signin/check-email", mustNotSignedIn(http.HandlerFunc(checkEmail)))
	main.Handle("/signin/link", mustNotSignedIn(http.HandlerFunc(signInLink)))
	main.Handle("/openid", mustNotSignedIn(http.HandlerFunc(openID)))
	main.Handle("/openid/callback", mustNotSignedIn(http.HandlerFunc(openIDCallback)))
	main.Handle("/signup", mustNotSignedIn(http.HandlerFunc(signUp)))
	main.Handle("/signout", http.HandlerFunc(signOut))
	main.Handle("/reset/password", mustNotSignedIn(http.HandlerFunc(resetPassword)))
	main.Handle("/profile", mustSignedIn(http.HandlerFunc(profile)))
	main.Handle("/profile/edit", mustSignedIn(http.HandlerFunc(profileEdit)))
	main.Handle("/course/", http.StripPrefix("/course/", courseHandler()))
	main.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
	main.Handle("/editor/", http.StripPrefix("/editor", editor))

	mux.Handle("/~/", http.StripPrefix("/~", cache(http.FileServer(&fileFS{http.Dir("static")}))))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

	mux.Handle("/", middleware.Chain(
		panicLogger,
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
		setDatabase(config.DB),
		fetchUser(),
		csrf(config.BaseURL, config.XSRFSecret),
	)(main))

	return middleware.Chain(
		setHeaders,
	)(mux)
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}
