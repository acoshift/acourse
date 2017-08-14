package app

import (
	"database/sql"
	"net/http"
	"os"
	"unicode/utf8"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/header"
	"github.com/acoshift/session"
	"github.com/asaskevich/govalidator"
)

// Handler returns app's handlers
func Handler() http.Handler {
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
	main.Handle("/openid", mustNotSignedIn(http.HandlerFunc(openID)))
	main.Handle("/openid/callback", mustNotSignedIn(http.HandlerFunc(openIDCallback)))
	main.Handle("/signup", mustNotSignedIn(http.HandlerFunc(signUp)))
	main.Handle("/signout", http.HandlerFunc(signOut))
	main.Handle("/profile", mustSignedIn(http.HandlerFunc(profile)))
	main.Handle("/profile/edit", mustSignedIn(http.HandlerFunc(profileEdit)))
	main.Handle("/course/", http.StripPrefix("/course/", http.HandlerFunc(course)))
	main.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin)))
	main.Handle("/editor/", http.StripPrefix("/editor", editor))
	main.Handle("/reset/password", mustNotSignedIn(http.HandlerFunc(resetPassword)))

	mux.Handle("/", Middleware(main))
	mux.Handle("/~/", http.StripPrefix("/~", cache(http.FileServer(&fileFS{http.Dir("static")}))))
	mux.Handle("/favicon.ico", fileHandler("static/favicon.ico"))

	return mux
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

func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.CacheControl, "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}

func fileHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		view.NotFound(w, r)
		return
	}
	courses, err := model.ListPublicCourses(ctx, db, cachePool, cachePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view.Index(w, r, courses)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postSignIn(w, r)
		return
	}
	view.SignIn(w, r)
}

func postSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := session.Get(ctx, sessName).Flash()

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

	s := session.Get(ctx, sessName)
	s.Set(keyUserID, userID)
	s.Rotate()

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		var id int64
		err = db.QueryRowContext(ctx, `select id from users where id = $1`, userID).Scan(&id)
		if err == sql.ErrNoRows {
			db.ExecContext(ctx, `insert into users (id, username, name, email) values ($1, $2, $3, $4)`, userID, userID, "", email)
		}
	}

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

func openID(w http.ResponseWriter, r *http.Request) {
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
	s := session.Get(ctx, sessName)
	s.Set(keyOpenIDSessionID, sessID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func openIDCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := session.Get(ctx, sessName)
	sessID, _ := s.Get(keyOpenIDSessionID).(string)
	s.Del(keyOpenIDSessionID)
	user, err := firAuth.VerifyAuthCallbackURI(ctx, baseURL+r.RequestURI, sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := db.BeginTx(ctx, nil)
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
	s.Rotate()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postSignUp(w, r)
		return
	}
	view.SignUp(w, r)
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := session.Get(ctx, sessName).Flash()

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

	_, err = db.ExecContext(ctx, `
		insert into users
			(id, username, name, email)
		values
			($1, $2, '', $3)
	`, userID, userID, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := session.Get(ctx, sessName)
	s.Set(keyUserID, userID)

	rURL := r.FormValue("r")
	if len(rURL) == 0 {
		rURL = "/"
	}

	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context(), sessName)
	s.Destroy()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer back(w, r)
		ctx := r.Context()
		f := session.Get(ctx, sessName).Flash()
		f.Set("OK", "1")
		email := r.FormValue("email")
		user, err := firAuth.GetUserByEmail(ctx, email)
		if err != nil {
			// don't send any error back to user
			return
		}
		firAuth.SendPasswordResetEmail(ctx, user.Email)
		return
	}
	view.ResetPassword(w, r)
}
