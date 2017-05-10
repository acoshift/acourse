package app

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
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
	gob.Register(sessionKey(0))

	mux := http.NewServeMux()

	mux.Handle("/", wrapFunc(getIndex, nil))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))
	mux.Handle("/signin", mustNotSignedIn(wrapFunc(getSignIn, postSignIn)))
	mux.Handle("/signup", mustNotSignedIn(wrapFunc(getSignUp, postSignUp)))
	mux.Handle("/signout", mustSignedIn(wrapFunc(getSignOut, nil)))
	mux.Handle("/profile", mustSignedIn(wrapFunc(getProfile, nil)))
	mux.Handle("/profile/edit", mustSignedIn(wrapFunc(getProfileEdit, postProfileEdit)))
	mux.Handle("/user/", http.StripPrefix("/user", wrapFunc(getUser, nil)))
	mux.Handle("/course/", http.StripPrefix("/course", wrapFunc(getCourse, nil)))

	admin := http.NewServeMux()
	// TODO: add admin route
	mux.Handle("/admin", onlyAdmin(admin))

	Handler = middleware.Chain(
		recovery,
		gzip.New(gzip.Config{Level: gzip.DefaultCompression}),
		session.Middleware(session.Config{
			Name:     "sess",
			Entropy:  32,
			Path:     "/",
			MaxAge:   10 * 24 * time.Hour,
			HTTPOnly: true,
			Secure:   session.PreferSecure,
			Store:    sRedis.New(internal.GetSecondaryPool(), "acr::s:"),
		}),
		flash.Middleware(),
		fetchUser,
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
	c := internal.GetPrimaryDB()
	defer c.Close()
	courses, err := model.ListPublicCourses(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Index(w, r, &view.IndexData{
		Page:    &defaultPage,
		Courses: courses,
	})
}

func getSignIn(w http.ResponseWriter, r *http.Request) {
	view.SignIn(w, r, &view.AuthData{
		Page:  &defaultPage,
		Flash: flash.Get(r.Context()),
	})
}

func postSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := flash.Get(ctx)

	if !verifyXSRF(r.FormValue("X"), "", "signin") {
		f.Add("Errors", "invalid xsrf token")
		back(w, r)
		return
	}

	email := r.FormValue("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}
	pass := r.FormValue("Password")
	if len(pass) == 0 {
		f.Add("Errors", "password required")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		back(w, r)
		return
	}

	// c := internal.GetPrimaryDB()
	// defer c.Close()

	userID, err := internal.SignInUser(email, pass)
	// u, err := model.GetUserFromEmailOrUsername(c, user)
	// if err == model.ErrNotFound {
	// 	f.Add("Errors", "wrong email or password")
	// 	back(w, r)
	// 	return
	// }
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	// if !verifyPassword(u.Password, pass) {
	// 	f.Add("Errors", "wrong email or password")
	// 	back(w, r)
	// 	return
	// }

	s := session.Get(ctx)
	// s.Set(keyUserID, u.ID())
	s.Set(keyUserID, userID)

	rURL := r.FormValue("r")
	if len(rURL) == 0 {
		rURL = "/"
	}

	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func getSignUp(w http.ResponseWriter, r *http.Request) {
	view.SignUp(w, r, &view.AuthData{
		Page: &defaultPage,
	})
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	defer back(w, r)
}

func getSignOut(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context())
	s.Del(keyUserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
