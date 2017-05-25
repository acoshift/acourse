package app

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/flash"
	"github.com/acoshift/header"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/session"
)

// Mount mounts app's handlers into mux
func Mount(mux *http.ServeMux) {
	r := httprouter.New()
	r.GET("/", http.HandlerFunc(getIndex))
	r.ServeFiles("/~/*filepath", http.Dir("static"))
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
	r.GET("/create-course", http.HandlerFunc(getCourseCreate))
	r.POST("/create-course", http.HandlerFunc(postCourseCreate))
	r.GET("/course/:courseID", http.HandlerFunc(getCourse))
	r.GET("/course/:courseID/edit", http.HandlerFunc(getCourseEdit))
	r.POST("/course/:courseID/edit", http.HandlerFunc(postCourseEdit))

	admin := httprouter.New()
	admin.GET("/users", http.HandlerFunc(getAdminUsers))
	admin.GET("/courses", http.HandlerFunc(getAdminCourses))
	admin.GET("/payments/pending", http.HandlerFunc(getAdminPendingPayments))
	admin.GET("/payments/history", http.HandlerFunc(getAdminHistoryPayments))

	mux.Handle("/", r)
	mux.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
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
	courses, err := model.ListPublicCourses()
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

	userID, err := SignInUser(email, pass)
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
	redirectURL, sessID, err := SignInUserProvider(p)
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
	userID, err := SignInUserProviderCallback(r.RequestURI, sessID)
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
	user := appctx.GetUser(r.Context())

	ownCourses, err := model.ListOwnCourses(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enrolledCourses, err := model.ListEnrolledCourses(user.ID)
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
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	f := flash.Get(ctx)
	if !f.Has("Username") {
		f.Set("Username", user.Username)
	}
	if !f.Has("Name") {
		f.Set("Name", user.Name)
	}
	if !f.Has("AboutMe") {
		f.Set("AboutMe", user.AboutMe)
	}
	page := defaultPage
	page.Title = user.Username + " | " + page.Title
	view.ProfileEdit(w, r, &view.ProfileEditData{
		Page:  &page,
		Flash: f,
	})
}

func postProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	f := flash.Get(ctx)
	if !verifyXSRF(r.FormValue("X"), user.ID, "profile-edit") {
		f.Add("Errors", "invalid xsrf token")
		back(w, r)
		return
	}
	image, info, err := r.FormFile("Image")
	var imageURL string
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	username := r.FormValue("Username")
	name := r.FormValue("Name")
	aboutMe := r.FormValue("AboutMe")
	user.Username = username
	user.Name = name
	user.AboutMe = aboutMe
	if len(imageURL) > 0 {
		user.Image = imageURL
	}
	f.Set("Username", username)
	f.Set("Name", name)
	f.Set("AboutMe", aboutMe)
	err = user.Save()
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	id := httprouter.GetParam(ctx, "courseID")
	course, err := model.GetCourFromIDOrURL(id)
	if err == model.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enrolled := false
	if user != nil {
		enrolled, err = model.IsEnrolled(user.ID, course.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	page := defaultPage
	page.Title = course.Title + " | " + page.Title
	page.Desc = course.ShortDesc
	page.Image = course.Image
	page.URL = baseURL + "/course/" + url.PathEscape(course.Link())
	view.Course(w, r, &view.CourseData{
		Page:     &page,
		Course:   course,
		Enrolled: enrolled,
	})
}

func getCourseCreate(w http.ResponseWriter, r *http.Request) {
	page := defaultPage
	page.Title = "Create new Course | " + page.Title
	view.CourseCreate(w, r, &view.CourseCreateData{
		Page: &page,
	})
}

func postCourseCreate(w http.ResponseWriter, r *http.Request) {

}

func getCourseEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := defaultPage
	page.Title = "Edit Course | " + page.Title

	courseID := httprouter.GetParam(ctx, "courseID")
	course, err := model.GetCourFromIDOrURL(courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.CourseEdit(w, r, &view.CourseEditData{
		Page:   &page,
		Course: course,
	})
}

func postCourseEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := httprouter.GetParam(ctx, "courseID")
	http.Redirect(w, r, "/course/"+id, http.StatusSeeOther)
}

func getAdminUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.AdminUsers(w, r, &view.AdminUsersData{
		Page:  &defaultPage,
		Users: users,
	})
}

func getAdminCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := model.ListCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.AdminCourses(w, r, &view.AdminCoursesData{
		Page:    &defaultPage,
		Courses: courses,
	})
}

func getAdminPayments(w http.ResponseWriter, r *http.Request, paymentsGetter func() ([]*model.Payment, error)) {
	action := r.FormValue("action")
	if len(action) > 0 {
		defer http.Redirect(w, r, "/admin/payments", http.StatusSeeOther)
		user := appctx.GetUser(r.Context())
		id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if err != nil {
			return
		}
		if action == "accept" && verifyXSRF(r.FormValue("x"), user.ID, "payment-accept") {
			x, err := model.GetPayment(id)
			if err != nil {
				return
			}
			x.Accept()
		} else if action == "reject" && verifyXSRF(r.FormValue("x"), user.ID, "payment-reject") {
			x, err := model.GetPayment(id)
			if err != nil {
				return
			}
			x.Reject()
		}
		return
	}

	payments, err := paymentsGetter()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.AdminPayments(w, r, &view.AdminPaymentsData{
		Page:     &defaultPage,
		Payments: payments,
	})
}

func getAdminPendingPayments(w http.ResponseWriter, r *http.Request) {
	getAdminPayments(w, r, model.ListPendingPayments)
}

func getAdminHistoryPayments(w http.ResponseWriter, r *http.Request) {
	getAdminPayments(w, r, model.ListHistoryPayments)
}
