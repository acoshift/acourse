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
	mux.Handle("/openid", mustNotSignedIn(wrapFunc(getSignInProvider, nil)))
	mux.Handle("/openid/callback", mustNotSignedIn(wrapFunc(getSignInCallback, nil)))
	mux.Handle("/signup", mustNotSignedIn(wrapFunc(getSignUp, postSignUp)))
	mux.Handle("/signout", mustSignedIn(wrapFunc(getSignOut, nil)))
	mux.Handle("/profile", mustSignedIn(wrapFunc(getProfile, nil)))
	mux.Handle("/profile/edit", mustSignedIn(wrapFunc(getProfileEdit, postProfileEdit)))
	// mux.Handle("/user/", http.StripPrefix("/user", wrapFunc(getUser, nil)))
	mux.Handle("/course/", http.StripPrefix("/course", &courseCtrl{}))

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

	userID, err := internal.SignInUser(email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	s := session.Get(ctx)
	s.Set(keyUserID, userID)

	rURL := r.FormValue("r")
	if len(rURL) == 0 {
		rURL = "/"
	}

	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

var allowProvider = map[string]bool{
	"google.com":   true,
	"facebook.com": true,
	"github.com":   true,
}

func getSignInProvider(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	if !allowProvider[p] {
		http.Error(w, "provider not allowed", http.StatusBadRequest)
		return
	}
	redirectURL, sessID, err := internal.SignInUserProvider(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := session.Get(r.Context())
	s.Set(keyOpenIDSessionID, sessID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func getSignInCallback(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context())
	sessID, _ := s.Get(keyOpenIDSessionID).(string)
	s.Del(keyOpenIDSessionID)
	userID, err := internal.SignInUserProviderCallback(r.RequestURI, sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Set(keyUserID, userID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getSignUp(w http.ResponseWriter, r *http.Request) {
	view.SignUp(w, r, &view.AuthData{
		Page:  &defaultPage,
		Flash: flash.Get(r.Context()),
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
	user, _ := internal.GetUser(r.Context()).(*model.User)

	c := internal.GetPrimaryDB()
	defer c.Close()
	ownCourses, err := model.ListOwnCourses(c, user.ID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enrolledCourses, err := model.ListEnrolledCourses(c, user.ID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := defaultPage
	page.Title = user.Username + " | " + page.Title

	view.Profile(w, r, &view.ProfileData{
		Page:            &page,
		OwnCourses:      ownCourses,
		EnrolledCourses: enrolledCourses,
	})
}

func getProfileEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "profile edit")
}

func postProfileEdit(w http.ResponseWriter, r *http.Request) {

}

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	ps := extractURL(r.URL)
// 	id := ps[0]
// 	if len(id) == 0 {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	fmt.Fprint(w, "user: ", id)
// }

type courseCtrl struct{}

func (*courseCtrl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps := extractURL(r.URL)
	if len(ps) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if len(ps) == 1 {
		ps = append(ps, "")
	}
	if len(ps) > 2 && len(ps[2]) == 0 {
		ps = ps[:2]
	}
	if len(ps) > 2 {
		http.NotFound(w, r)
		return
	}
	id := ps[0]
	if len(id) == 0 {
		http.NotFound(w, r)
		return
	}
	switch ps[1] {
	case "":
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			getCourse(w, r, id)
			return
		}
	case "edit":
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			getCourseEdit(w, r, id)
			return
		}
		if r.Method == http.MethodPost {
			postCourseEdit(w, r, id)
			return
		}
	default:
		http.NotFound(w, r)
		return
	}
	status := http.StatusMethodNotAllowed
	http.Error(w, http.StatusText(status), status)
}

func getCourse(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Fprint(w, "course: ", id)
}

func getCourseEdit(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Fprint(w, "course edit: ", id)
}

func postCourseEdit(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Fprint(w, "course edit: ", id)
}
