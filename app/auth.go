package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/acoshift/go-firebase-admin"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func generateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateSessionID() string {
	return generateRandomString(24)
}

func generateMagicLinkID() string {
	return generateRandomString(64)
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
	s := appctx.GetSession(ctx)
	f := s.Flash()

	email := r.FormValue("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		back(w, r)
		return
	}

	ok, err := repository.CanAcquireMagicLink(ctx, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		f.Add("Errors", "อีเมลของคุณได้ขอ Magic Link จากเราไปแล้ว กรุณาตรวจสอบอีเมล")
		back(w, r)
	}

	f.Set("CheckEmail", "1")

	user, err := repository.FindUserByEmail(ctx, email)
	// don't lets user know if email is wrong
	if err == appctx.ErrNotFound {
		http.Redirect(w, r, "/signin/check-email", http.StatusSeeOther)
		return
	}
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	linkID := generateMagicLinkID()

	err = repository.StoreMagicLink(ctx, linkID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	linkQuery := make(url.Values)
	linkQuery.Set("id", linkID)
	if x := r.FormValue("r"); len(x) > 0 {
		linkQuery.Set("r", parsePath(x))
	}

	message := fmt.Sprintf(`สวัสดีครับคุณ %s,


ตามที่ท่านได้ขอ Magic Link เพื่อเข้าสู่ระบบสำหรับ acourse.io นั้นท่านสามารถเข้าได้ผ่าน Link ข้างล่างนี้ ภายใน 1 ชม.

%s

ทีมงาน acourse.io
	`, user.Name, makeLink("/signin/link", linkQuery))

	go sendEmail(user.Email.String, "Magic Link Request", markdown(message))

	http.Redirect(w, r, "/signin/check-email", http.StatusSeeOther)
}

func checkEmail(w http.ResponseWriter, r *http.Request) {
	f := appctx.GetSession(r.Context()).Flash()
	if !f.Has("CheckEmail") {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	view.CheckEmail(w, r)
}

func signInLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	linkID := r.FormValue("id")
	if len(linkID) == 0 {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	s := appctx.GetSession(ctx)
	f := s.Flash()

	userID, err := repository.FindMagicLink(ctx, linkID)
	if err != nil {
		f.Add("Errors", "ไม่พบ Magic Link ของคุณ")
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setUserID(s, userID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func signInPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postSignInPassword(w, r)
		return
	}
	view.SignInPassword(w, r)
}

func postSignInPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := appctx.GetSession(ctx)
	f := s.Flash()

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

	userID, err := auth.VerifyPassword(ctx, email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	s.Rotate()
	setUserID(s, userID)

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		ok, err := repository.IsUserExists(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			err = repository.CreateUser(ctx, &entity.User{ID: userID, Email: sql.NullString{String: email, Valid: len(email) > 0}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	rURL := parsePath(r.FormValue("r"))
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
	redirectURL, err := auth.CreateAuthURI(ctx, p, baseURL+"/openid/callback", sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s := appctx.GetSession(ctx)
	setOpenIDSessionID(s, sessID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func openIDCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := appctx.GetSession(ctx)
	sessID := getOpenIDSessionID(s)
	delOpenIDSessionID(s)
	user, err := auth.VerifyAuthCallbackURI(ctx, baseURL+r.RequestURI, sessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, tx, err := appctx.NewTransactionContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	db := appctx.GetTransaction(ctx)
	// check is user sign up
	var cnt int64
	err = db.QueryRowContext(ctx, `select 1 from users where id = $1`, user.UserID).Scan(&cnt)
	if err == sql.ErrNoRows {
		// user not found, insert new user
		imageURL := uploadProfileFromURLAsync(user.PhotoURL)
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

	s.Rotate()
	setUserID(s, user.UserID)
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
	f := appctx.GetSession(ctx).Flash()

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

	userID, err := auth.CreateUser(ctx, &firebase.User{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	db := appctx.GetDatabase(ctx)
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

	s := appctx.GetSession(ctx)
	setUserID(s, userID)

	rURL := parsePath(r.FormValue("r"))
	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	appctx.GetSession(r.Context()).Destroy()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer back(w, r)
		ctx := r.Context()
		f := appctx.GetSession(ctx).Flash()
		f.Set("OK", "1")
		email := r.FormValue("email")
		user, err := auth.GetUserByEmail(ctx, email)
		if err != nil {
			// don't send any error back to user
			return
		}
		auth.SendPasswordResetEmail(ctx, user.Email)
		return
	}
	view.ResetPassword(w, r)
}
