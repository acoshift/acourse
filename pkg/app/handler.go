package app

import (
	"database/sql"
	"net/http"
	"os"
	"unicode/utf8"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/flash"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	"github.com/asaskevich/govalidator"
)

// Mount mounts app's handlers into mux
func Mount(mux *http.ServeMux) {
	r := httprouter.New()
	r.GET("/", http.HandlerFunc(getIndex))

	r.GET("/signin", mustNotSignedIn(http.HandlerFunc(getSignIn)))
	r.POST("/signin", middleware.Chain(
		mustNotSignedIn,
		xsrf("signin"),
	)(http.HandlerFunc(postSignIn)))
	r.GET("/openid", mustNotSignedIn(http.HandlerFunc(getSignInProvider)))
	r.GET("/openid/callback", mustNotSignedIn(http.HandlerFunc(getSignInCallback)))
	r.GET("/signup", mustNotSignedIn(http.HandlerFunc(getSignUp)))
	r.POST("/signup", middleware.Chain(
		mustNotSignedIn,
		xsrf("signup"),
	)(http.HandlerFunc(postSignUp)))
	r.GET("/signout", http.HandlerFunc(getSignOut))

	r.GET("/profile", mustSignedIn(http.HandlerFunc(getProfile)))
	r.GET("/profile/edit", mustSignedIn(http.HandlerFunc(getProfileEdit)))
	r.POST("/profile/edit", middleware.Chain(
		mustSignedIn,
		xsrf("profile/edit"),
	)(http.HandlerFunc(postProfileEdit)))

	r.GET("/course/:courseID", http.HandlerFunc(getCourse))
	r.GET("/course/:courseID/content", mustSignedIn(http.HandlerFunc(getCourseContent)))
	r.GET("/course/:courseID/enroll", mustSignedIn(http.HandlerFunc(getCourseEnroll)))
	r.POST("/course/:courseID/enroll", middleware.Chain(
		mustSignedIn,
		xsrf("enroll"),
	)(http.HandlerFunc(postCourseEnroll)))

	editor := http.NewServeMux()
	{
		post := xsrf("editor/create")(http.HandlerFunc(postEditorCreate))
		editor.Handle("/create", onlyInstructor(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead:
				getEditorCreate(w, r)
			case http.MethodPost:
				post.ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
			}
		})))
	}
	{
		post := xsrf("editor/course")(http.HandlerFunc(postEditorCourse))
		editor.Handle("/course", isCourseOwner(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead:
				getEditorCourse(w, r)
			case http.MethodPost:
				post.ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
			}
		})))
	}
	editor.Handle("/content", isCourseOwner(http.HandlerFunc(getEditorContent)))
	editor.Handle("/content/create", isCourseOwner(http.HandlerFunc(getEditorContentCreate)))
	editor.Handle("/content/edit", http.HandlerFunc(getEditorContentEdit))

	admin := http.NewServeMux()
	admin.Handle("/users", http.HandlerFunc(getAdminUsers))
	admin.Handle("/courses", http.HandlerFunc(getAdminCourses))
	{
		post := xsrf("payment-action")(http.HandlerFunc(postAdminPendingPayment))
		admin.Handle("/payments/pending", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead:
				getAdminPendingPayments(w, r)
			case http.MethodPost:
				post.ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
			}
		}))
	}
	admin.Handle("/payments/history", http.HandlerFunc(getAdminHistoryPayments))

	mux.Handle("/", r)
	mux.Handle("/~/", http.StripPrefix("/~", http.FileServer(&fileFS{http.Dir("static")})))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))
	mux.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
	mux.Handle("/editor/", http.StripPrefix("/editor", editor))
}

type fileFS struct {
	http.FileSystem
}

func (fs *fileFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
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
	view.Index(w, r, courses)
}

func getSignIn(w http.ResponseWriter, r *http.Request) {
	view.SignIn(w, r)
}

func postSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := flash.Get(ctx)

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

	userID, err := firAuth.VerifyPassword(ctx, email, pass)
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

	ctx := r.Context()
	sessID := generateSessionID()
	redirectURL, err := firAuth.CreateAuthURI(ctx, p, baseURL+"/openid/callback", sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := session.Get(ctx)
	s.Set(keyOpenIDSessionID, sessID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func getSignInCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := session.Get(ctx)
	sessID, _ := s.Get(keyOpenIDSessionID).(string)
	s.Del(keyOpenIDSessionID)
	user, err := firAuth.VerifyAuthCallbackURI(ctx, baseURL+r.RequestURI, sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// check is user sign up
	var cnt int64
	err = tx.QueryRow(`select 1 from users where id = $1`, user.UserID).Scan(&cnt)
	if err == sql.ErrNoRows {
		// user not found, insert new user
		imageURL := UploadProfileFromURLAsync(user.PhotoURL)
		_, err = tx.Exec(`
			insert into users
				(id, name, username, email, image)
			values
				($1, $2, $3, $4, $5)
		`, user.UserID, user.DisplayName, user.UserID, sql.NullString{String: user.Email, Valid: len(user.Email) > 0}, imageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Set(keyUserID, user.UserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getSignUp(w http.ResponseWriter, r *http.Request) {
	view.SignUp(w, r)
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := flash.Get(ctx)

	email := r.FormValue("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		f.Add("Errors", err.Error())
		return
	}
	pass := r.FormValue("Password")
	if len(pass) == 0 {
		f.Add("Errors", "password required")
	}
	if n := utf8.RuneCountInString(pass); n < 6 || n > 64 {
		f.Add("Errors", "password must have 6 to 64 characters")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		back(w, r)
		return
	}

	userID, err := firAuth.CreateUser(ctx, &admin.User{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	_, err = db.Exec(`
		insert into users
			(id, username, email)
		values
			($1, $2, $3)
	`, userID, userID, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func getSignOut(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context())
	s.Del(keyUserID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
