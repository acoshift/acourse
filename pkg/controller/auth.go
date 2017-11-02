package controller

import (
	"database/sql"
	"net/http"
	"unicode/utf8"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/session"
	"github.com/asaskevich/govalidator"
)

func (c *ctrl) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postSignIn(w, r)
		return
	}
	c.view.SignIn(w, r)
}

func (c *ctrl) postSignIn(w http.ResponseWriter, r *http.Request) {
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

	userID, err := c.firAuth.VerifyPassword(ctx, email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	s := session.Get(ctx, sessName)
	app.SetUserID(s, userID)
	s.Rotate()

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		var id string
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

func (c *ctrl) openID(w http.ResponseWriter, r *http.Request) {
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

	app.SetUserID(s, user.UserID)
	s.Rotate()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *ctrl) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postSignUp(w, r)
		return
	}
	c.view.SignUp(w, r)
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

	userID, err := firAuth.CreateUser(ctx, &firebase.User{
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
	app.SetUserID(s, userID)

	rURL := r.FormValue("r")
	if len(rURL) == 0 {
		rURL = "/"
	}

	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func (c *ctrl) signOut(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context(), sessName)
	s.Destroy()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *ctrl) resetPassword(w http.ResponseWriter, r *http.Request) {
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
	c.view.ResetPassword(w, r)
}
