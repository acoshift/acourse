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
	"github.com/acoshift/hime"
	"github.com/acoshift/pgsql"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
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

func signIn(ctx hime.Context) hime.Result {
	return ctx.View("signin", newPage(ctx))
}

func postSignIn(ctx hime.Context) hime.Result {
	s := appctx.GetSession(ctx)
	f := s.Flash()

	email := ctx.FormValueTrimSpace("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}
	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		f.Add("Errors", "invalid email")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		return ctx.RedirectToGet()
	}

	ok, err := repository.CanAcquireMagicLink(redisPool, redisPrefix, email)
	must(err)
	if !ok {
		f.Add("Errors", "อีเมลของคุณได้ขอ Magic Link จากเราไปแล้ว กรุณาตรวจสอบอีเมล")
		return ctx.RedirectToGet()
	}

	f.Set("CheckEmail", "1")

	user, err := repository.FindUserByEmail(db, email)
	// don't lets user know if email is wrong
	if err == entity.ErrNotFound {
		return ctx.RedirectTo("signin.check-email")
	}
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	linkID := generateMagicLinkID()

	err = repository.StoreMagicLink(redisPool, redisPrefix, linkID, user.ID)
	must(err)

	linkQuery := make(url.Values)
	linkQuery.Set("id", linkID)

	message := fmt.Sprintf(`สวัสดีครับคุณ %s,


ตามที่ท่านได้ขอ Magic Link เพื่อเข้าสู่ระบบสำหรับ acourse.io นั้นท่านสามารถเข้าได้ผ่าน Link ข้างล่างนี้ ภายใน 1 ชม.

%s

ทีมงาน acourse.io
	`, user.Name, baseURL+"/signin/link?id="+linkID)

	go sendEmail(user.Email.String, "Magic Link Request", markdown(message))

	return ctx.RedirectTo("signin.check-email")
}

func checkEmail(ctx hime.Context) hime.Result {
	f := appctx.GetSession(ctx).Flash()
	if !f.Has("CheckEmail") {
		return ctx.Redirect("/")
	}
	return ctx.View("check-email", newPage(ctx))
}

func signInLink(ctx hime.Context) hime.Result {
	linkID := ctx.FormValue("id")
	if len(linkID) == 0 {
		return ctx.RedirectTo("signin")
	}

	s := appctx.GetSession(ctx)
	f := s.Flash()

	userID, err := repository.FindMagicLink(redisPool, redisPrefix, linkID)
	if err != nil {
		f.Add("Errors", "ไม่พบ Magic Link ของคุณ")
		return ctx.RedirectTo("signin")
	}

	setUserID(s, userID)
	return ctx.Redirect("/")
}

func signInPassword(ctx hime.Context) hime.Result {
	return ctx.View("signin.password", newPage(ctx))
}

func postSignInPassword(ctx hime.Context) hime.Result {
	s := appctx.GetSession(ctx)
	f := s.Flash()

	email := ctx.PostFormValue("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}
	pass := ctx.PostFormValue("Password")
	if len(pass) == 0 {
		f.Add("Errors", "password required")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		return ctx.RedirectToGet()
	}

	userID, err := auth.VerifyPassword(ctx, email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	s.Regenerate()
	setUserID(s, userID)

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		ok, err := repository.IsUserExists(db, userID)
		must(err)
		if !ok {
			err = repository.CreateUser(db, &entity.User{ID: userID, Email: sql.NullString{String: email, Valid: len(email) > 0}})
			must(err)
		}
	}

	return ctx.SafeRedirect(ctx.FormValue("r"))
}

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func openID(ctx hime.Context) hime.Result {
	p := ctx.FormValue("p")
	if !allowProvider[p] {
		return ctx.Status(http.StatusBadRequest).String("provider not allowed")
	}

	sessID := generateSessionID()
	redirectURL, err := auth.CreateAuthURI(ctx, p, baseURL+"/openid/callback", sessID)
	must(err)

	s := appctx.GetSession(ctx)
	setOpenIDSessionID(s, sessID)
	return ctx.Redirect(redirectURL)
}

func openIDCallback(ctx hime.Context) hime.Result {
	s := appctx.GetSession(ctx)
	sessID := getOpenIDSessionID(s)
	delOpenIDSessionID(s)
	user, err := auth.VerifyAuthCallbackURI(ctx, baseURL+ctx.Request().RequestURI, sessID)
	must(err)

	err = pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		// check is user sign up
		var cnt int64
		err = tx.QueryRowContext(ctx, `select 1 from users where id = $1`, user.UserID).Scan(&cnt)
		if err == sql.ErrNoRows {
			// user not found, insert new user
			imageURL := uploadProfileFromURLAsync(user.PhotoURL)
			_, err = tx.ExecContext(ctx, `
				insert into users
					(id, name, username, email, image)
				values
					($1, $2, $3, $4, $5)
			`, user.UserID, user.DisplayName, user.UserID, sql.NullString{String: user.Email, Valid: len(user.Email) > 0}, imageURL)
			if err != nil {
				return err
			}
		}
		return nil
	})
	must(err)

	s.Regenerate()
	setUserID(s, user.UserID)
	return ctx.RedirectTo("index")
}

func signUp(ctx hime.Context) hime.Result {
	return ctx.View("signup", newPage(ctx))
}

func postSignUp(ctx hime.Context) hime.Result {
	f := appctx.GetSession(ctx).Flash()

	email := ctx.FormValue("Email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	pass := ctx.FormValue("Password")
	if len(pass) == 0 {
		f.Add("Errors", "password required")
	}
	if n := utf8.RuneCountInString(pass); n < 6 || n > 64 {
		f.Add("Errors", "password must have 6 to 64 characters")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		return ctx.RedirectToGet()
	}

	userID, err := auth.CreateUser(ctx, &firebase.User{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	_, err = db.ExecContext(ctx, `
		insert into users
			(id, username, name, email)
		values
			($1, $2, '', $3)
	`, userID, userID, email)
	must(err)

	s := appctx.GetSession(ctx)
	setUserID(s, userID)

	return ctx.SafeRedirect(ctx.FormValue("r"))
}

func signOut(ctx hime.Context) hime.Result {
	appctx.GetSession(ctx).Destroy()
	return ctx.Redirect("/")
}

func resetPassword(ctx hime.Context) hime.Result {
	return ctx.View("reset.password", newPage(ctx))
}

func postResetPassword(ctx hime.Context) hime.Result {
	f := appctx.GetSession(ctx).Flash()
	f.Set("OK", "1")
	email := ctx.FormValue("email")
	user, err := auth.GetUserByEmail(ctx, email)
	if err != nil {
		// don't send any error back to user
		return ctx.RedirectToGet()
	}
	auth.SendPasswordResetEmail(ctx, user.Email)

	return ctx.RedirectToGet()
}
