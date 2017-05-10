package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/view"
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
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))
	mux.Handle("/signin", mustNotSignedIn(wrapFunc(getSignIn, postSignIn)))
	mux.Handle("/signout", mustSignedIn(wrapFunc(getSignOut, nil)))
	mux.Handle("/profile", mustSignedIn(wrapFunc(getProfile, nil)))
	mux.Handle("/profile/edit", mustSignedIn(wrapFunc(getProfileEdit, postProfileEdit)))
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

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}

var defaultPage = view.Page{
	Title: "Acourse",
	Desc:  "Online courses for everyone",
	Image: "https://storage.googleapis.com/acourse/static/62b9eb0e-3668-4f9f-86b7-a11349938f7a.jpg",
	URL:   "https://acourse.io",
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	view.Index(w, r, &view.IndexData{
		Page: &defaultPage,
	})
}

func getSignIn(w http.ResponseWriter, r *http.Request) {
	view.SignIn(w, r, &view.AuthData{
		Page: &defaultPage,
	})
}

func postSignIn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusFound)
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
	id := extractPathID(r.URL)
	if len(id) == 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "user: ", id)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	id := extractPathID(r.URL)
	if len(id) == 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "course: ", id)
}
