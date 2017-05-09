package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/flash"
	"github.com/acoshift/gzip"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	sRedis "github.com/acoshift/session/store/redis"
)

// Handler returns app's handler
var Handler http.Handler

func init() {
	mux := http.NewServeMux()

	mux.Handle("/", wrapFunc(getIndex, nil))
	mux.Handle("/signin", wrapFunc(getSignIn, postSignIn))
	mux.Handle("/signout", wrapFunc(getSignOut, nil))
	mux.Handle("/profile", wrapFunc(getProfile, nil))
	mux.Handle("/profile/edit", wrapFunc(getProfileEdit, postProfileEdit))
	mux.Handle("/user/", http.StripPrefix("/user", wrapFunc(getUser, nil)))
	mux.Handle("/course/", http.StripPrefix("/course", wrapFunc(getCourse, nil)))

	Handler = middleware.Chain(
		recovery,
		gzip.New(gzip.Config{Level: gzip.DefaultCompression}),
		session.Middleware(session.Config{
			Name:     "sess",
			Path:     "/",
			MaxAge:   10 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			Store:    sRedis.New(internal.GetSecondaryPool(), "acr:s:"),
		}),
		flash.Middleware(),
	)(mux)
}

func wrapFunc(get, post http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			if get != nil {
				get.ServeHTTP(w, r)
				return
			}
		case http.MethodPost:
			if post != nil {
				post.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "index")
}

func getSignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "signIn")
}

func postSignIn(w http.ResponseWriter, r *http.Request) {

}

func getSignOut(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context())
	s.Del(keyUserID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "profile")
}

func getProfileEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "profile edit")
}

func postProfileEdit(w http.ResponseWriter, r *http.Request) {

}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "user", r.URL)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "course", r.URL)
}
