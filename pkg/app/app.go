package app

import (
	"net/http"
	"time"

	"github.com/acoshift/cachestatic"
	"github.com/acoshift/middleware"
	"github.com/acoshift/servertiming"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
)

// New creates new app
func New(config Config) http.Handler {
	ctrl := config.Controller
	repo := config.Repository
	view := config.View

	app := &app{
		ctrl: ctrl,
		repo: repo,
		view: view,
	}

	cacheInvalidator := make(chan interface{})

	go func() {
		for {
			time.Sleep(15 * time.Second)
			cacheInvalidator <- cachestatic.InvalidateAll
		}
	}()

	// create middlewares
	isCourseOwner := isCourseOwner(config.DB, config.View)

	// create mux
	mux := http.NewServeMux()

	editor := http.NewServeMux()
	editor.Handle("/create", onlyInstructor(http.HandlerFunc(ctrl.EditorCreate)))
	editor.Handle("/course", isCourseOwner(http.HandlerFunc(ctrl.EditorCourse)))
	editor.Handle("/content", isCourseOwner(http.HandlerFunc(ctrl.EditorContent)))
	editor.Handle("/content/create", isCourseOwner(http.HandlerFunc(ctrl.EditorContentCreate)))
	editor.Handle("/content/edit", http.HandlerFunc(ctrl.EditorContentEdit))

	admin := http.NewServeMux()
	admin.Handle("/users", http.HandlerFunc(ctrl.AdminUsers))
	admin.Handle("/courses", http.HandlerFunc(ctrl.AdminCourses))
	admin.Handle("/payments/pending", http.HandlerFunc(ctrl.AdminPendingPayments))
	admin.Handle("/payments/history", http.HandlerFunc(ctrl.AdminHistoryPayments))
	admin.Handle("/payments/reject", http.HandlerFunc(ctrl.AdminRejectPayment))

	main := http.NewServeMux()
	main.Handle("/", http.HandlerFunc(ctrl.Index))
	main.Handle("/signin", mustNotSignedIn(http.HandlerFunc(ctrl.SignIn)))
	main.Handle("/signin/password", mustNotSignedIn(http.HandlerFunc(ctrl.SignInPassword)))
	main.Handle("/signin/check-email", mustNotSignedIn(http.HandlerFunc(ctrl.CheckEmail)))
	main.Handle("/signin/link", mustNotSignedIn(http.HandlerFunc(ctrl.SignInLink)))
	main.Handle("/openid", mustNotSignedIn(http.HandlerFunc(ctrl.OpenID)))
	main.Handle("/openid/callback", mustNotSignedIn(http.HandlerFunc(ctrl.OpenIDCallback)))
	main.Handle("/signup", mustNotSignedIn(http.HandlerFunc(ctrl.SignUp)))
	main.Handle("/signout", http.HandlerFunc(ctrl.SignOut))
	main.Handle("/reset/password", mustNotSignedIn(http.HandlerFunc(ctrl.ResetPassword)))
	main.Handle("/profile", mustSignedIn(http.HandlerFunc(ctrl.Profile)))
	main.Handle("/profile/edit", mustSignedIn(http.HandlerFunc(ctrl.ProfileEdit)))
	main.Handle("/course/", http.StripPrefix("/course/", courseHandler(ctrl, view)))
	main.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
	main.Handle("/editor/", http.StripPrefix("/editor", editor))

	mux.Handle("/~/", http.StripPrefix("/~", cache(http.FileServer(&fileFS{http.Dir("static")}))))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

	mux.Handle("/", middleware.Chain(
		servertiming.Middleware(),
		panicLogger,
		session.Middleware(session.Config{
			Secret:   config.SessionSecret,
			Path:     "/",
			MaxAge:   5 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			Store: redisstore.New(redisstore.Config{
				Prefix: config.RedisPrefix,
				Pool:   config.RedisPool,
			}),
		}),
		cachestatic.New(cachestatic.Config{
			Skipper: func(r *http.Request) bool {
				// cache only get
				if r.Method != http.MethodGet {
					return true
				}

				// skip if signed in
				s := session.Get(r.Context(), sessName)
				if x := GetUserID(s); len(x) > 0 {
					return true
				}

				// cache only index
				if r.URL.Path == "/" {
					return false
				}
				return true
			},
			Invalidator: cacheInvalidator,
		}),
		setDatabase(config.DB),
		setRedisPool(config.RedisPool),
		fetchUser(repo),
		csrf(config.BaseURL, config.XSRFSecret),
	)(main))

	app.Handler = middleware.Chain(
		setHeaders,
	)(mux)

	return app
}

type app struct {
	http.Handler
	ctrl Controller
	repo Repository
	view View
}
