package controller

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/pkg/app"
)

func generateSessionID() string {
	b := make([]byte, 24)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (c *ctrl) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postSignIn(w, r)
		return
	}
	c.view.SignIn(w, r)
}

func (c *ctrl) postSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := app.GetSession(ctx).Flash()

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

	userID, err := c.auth.VerifyPassword(ctx, email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	s := app.GetSession(ctx)
	app.SetUserID(s, userID)
	s.Rotate()

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		ok, err := c.repo.IsUserExists(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			err = c.repo.CreateUser(ctx, &app.User{ID: userID, Email: sql.NullString{String: email, Valid: len(email) > 0}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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

func (c *ctrl) OpenID(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	if !allowProvider[p] {
		http.Error(w, "provider not allowed", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	sessID := generateSessionID()
	redirectURL, err := c.auth.CreateAuthURI(ctx, p, c.baseURL+"/openid/callback", sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := app.GetSession(ctx)
	app.SetOpenIDSessionID(s, sessID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (c *ctrl) OpenIDCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := app.GetSession(ctx)
	sessID := app.GetOpenIDSessionID(s)
	app.DelOpenIDSessionID(s)
	user, err := c.auth.VerifyAuthCallbackURI(ctx, c.baseURL+r.RequestURI, sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, tx, err := app.WithTransaction(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	db := app.GetTransaction(ctx)
	// check is user sign up
	var cnt int64
	err = db.QueryRowContext(ctx, `select 1 from users where id = $1`, user.UserID).Scan(&cnt)
	if err == sql.ErrNoRows {
		// user not found, insert new user
		imageURL := c.uploadProfileFromURLAsync(user.PhotoURL)
		_, err = db.ExecContext(ctx, `
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

func (c *ctrl) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postSignUp(w, r)
		return
	}
	c.view.SignUp(w, r)
}

func (c *ctrl) postSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f := app.GetSession(ctx).Flash()

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

	userID, err := c.auth.CreateUser(ctx, &firebase.User{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	db := app.GetDatabase(ctx)
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

	s := app.GetSession(ctx)
	app.SetUserID(s, userID)

	rURL := r.FormValue("r")
	if len(rURL) == 0 {
		rURL = "/"
	}

	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func (c *ctrl) SignOut(w http.ResponseWriter, r *http.Request) {
	app.GetSession(r.Context()).Destroy()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *ctrl) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer back(w, r)
		ctx := r.Context()
		f := app.GetSession(ctx).Flash()
		f.Set("OK", "1")
		email := r.FormValue("email")
		user, err := c.auth.GetUserByEmail(ctx, email)
		if err != nil {
			// don't send any error back to user
			return
		}
		c.auth.SendPasswordResetEmail(ctx, user.Email)
		return
	}
	c.view.ResetPassword(w, r)
}
