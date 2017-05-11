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
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	sRedis "github.com/acoshift/session/store/redis"
)

// Handler returns app's handler
var Handler http.Handler

func init() {
	gob.Register(sessionKey(0))

	mux := http.NewServeMux()

	r := httprouter.New()
	r.GET("/", http.HandlerFunc(getIndex))
	r.GET("/favicon.ico", fileHandler("static/favicon.ico"))
	r.GET("/signin", mustNotSignedIn(http.HandlerFunc(getSignIn)))
	r.POST("/signin", mustNotSignedIn(http.HandlerFunc(postSignIn)))
	r.GET("/openid", mustNotSignedIn(http.HandlerFunc(getSignInProvider)))
	r.GET("/openid/callback", mustNotSignedIn(http.HandlerFunc(getSignInCallback)))
	r.GET("/signup", mustNotSignedIn(http.HandlerFunc(getSignUp)))
	r.POST("/signup", mustNotSignedIn(http.HandlerFunc(postSignUp)))
	r.GET("/signout", mustSignedIn(http.HandlerFunc(getSignOut)))
	r.GET("/profile", mustSignedIn(http.HandlerFunc(getProfile)))
	r.GET("/profile/edit", mustSignedIn(http.HandlerFunc(getProfileEdit)))
	r.POST("/profile/edit", mustSignedIn(http.HandlerFunc(postProfileEdit)))
	r.GET("/course/:courseID", http.HandlerFunc(getCourse))
	r.GET("/course/:courseID/edit", http.HandlerFunc(getCourseEdit))
	r.POST("/course/:courseID/edit", http.HandlerFunc(postCourseEdit))

	admin := httprouter.New()
	admin.GET("/users", http.HandlerFunc(getAdminUsers))
	admin.GET("/courses", http.HandlerFunc(getAdminCourses))
	admin.GET("/payments", http.HandlerFunc(getAdminPayments))

	mux.Handle("/", r)
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
	c.Close()
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
	c.Close()

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

func getCourse(w http.ResponseWriter, r *http.Request) {
	id := httprouter.GetParam(r.Context(), "courseID")
	fmt.Fprint(w, "course: ", id)
}

func getCourseEdit(w http.ResponseWriter, r *http.Request) {
	id := httprouter.GetParam(r.Context(), "courseID")
	fmt.Fprint(w, "course edit: ", id)
}

func postCourseEdit(w http.ResponseWriter, r *http.Request) {
	id := httprouter.GetParam(r.Context(), "courseID")
	fmt.Fprint(w, "course edit: ", id)
}

func getAdminUsers(w http.ResponseWriter, r *http.Request) {
	c := internal.GetPrimaryDB()
	defer c.Close()
	users, err := model.ListUsers(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Close()
	view.AdminUsers(w, r, &view.AdminUsersData{
		Page:  &defaultPage,
		Users: users,
	})
}

func getAdminCourses(w http.ResponseWriter, r *http.Request) {
	view.AdminCourses(w, r, &view.AdminCoursesData{
		Page: &defaultPage,
	})
}

func getAdminPayments(w http.ResponseWriter, r *http.Request) {
	view.AdminPayments(w, r, &view.AdminPaymentsData{
		Page: &defaultPage,
	})
}
