package auth

import (
	"context"
	"net/http"

	"github.com/acoshift/hime"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

var allowProvider = map[string]bool{
	"google.com": true,
	"github.com": true,
}

func generateSessionID() string {
	return generateRandomString(24)
}

func (c *ctrl) openID(ctx *hime.Context) error {
	p := ctx.FormValue("p")
	if !allowProvider[p] {
		return ctx.Status(http.StatusBadRequest).String("provider not allowed")
	}

	sessID := generateSessionID()
	redirectURL, err := c.Auth.CreateAuthURI(ctx, p, c.BaseURL+ctx.Route("auth.openid.callback"), sessID)
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
