package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/hime"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
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

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func (c *ctrl) openID(ctx *hime.Context) error {
	p := ctx.FormValue("p")
	if !allowProvider[p] {
		return ctx.Status(http.StatusBadRequest).String("provider not allowed")
	}

	sessID := generateSessionID()
	redirectURL, err := c.Auth.CreateAuthURI(ctx, p, c.BaseURL+ctx.Route("openid.callback"), sessID)
	if err != nil {
		return err
	}

	s := appctx.GetSession(ctx)
	appsess.SetOpenIDSessionID(s, sessID)
	return ctx.Redirect(redirectURL)
}

func (c *ctrl) openIDCallback(ctx *hime.Context) error {
	s := appctx.GetSession(ctx)
	sessID := appsess.GetOpenIDSessionID(s)
	appsess.DelOpenIDSessionID(s)
	user, err := c.Auth.VerifyAuthCallbackURI(ctx, c.BaseURL+ctx.Request().RequestURI, sessID)
	if err != nil {
		return err
	}

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		// check is user sign up
		exists, err := repository.IsUserExists(ctx, user.UserID)
		if err != nil {
			return err
		}
		if !exists {
			// user not found, insert new user
			imageURL := c.uploadProfileFromURLAsync(user.PhotoURL)
			err = repository.RegisterUser(ctx, &entity.RegisterUser{
				ID:       user.UserID,
				Name:     user.DisplayName,
				Username: user.UserID,
				Email:    user.Email,
				Image:    imageURL,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.Regenerate()
	appsess.SetUserID(s, user.UserID)
	return ctx.RedirectTo("index")
}

func (c *ctrl) signUp(ctx *hime.Context) error {
	return ctx.View("signup", view.Page(ctx))
}

func (c *ctrl) postSignUp(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()

	email := ctx.FormValue("email")
	if len(email) == 0 {
		f.Add("Errors", "email required")
	}

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}
	pass := ctx.FormValue("password")
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

	userID, err := c.Auth.CreateUser(ctx, &firebase.User{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	err = repository.RegisterUser(ctx, &entity.RegisterUser{
		ID:       userID,
		Username: userID,
		Email:    email,
	})
	if err != nil {
		return err
	}

	s := appctx.GetSession(ctx)
	appsess.SetUserID(s, userID)

	return ctx.SafeRedirect(ctx.FormValue("r"))
}

func (c *ctrl) resetPassword(ctx *hime.Context) error {
	return ctx.View("reset.password", view.Page(ctx))
}

func (c *ctrl) postResetPassword(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()
	f.Set("OK", "1")
	email := ctx.FormValue("email")
	user, err := c.Auth.GetUserByEmail(ctx, email)
	if err != nil {
		// don't send any error back to user
		return ctx.RedirectToGet()
	}
	c.Auth.SendPasswordResetEmail(ctx, user.Email)

	return ctx.RedirectToGet()
}
