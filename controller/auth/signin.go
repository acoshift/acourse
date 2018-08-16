package auth

import (
	"fmt"
	"net/url"

	"github.com/acoshift/hime"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appsess"
	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) signIn(ctx *hime.Context) error {
	return ctx.View("auth.signin", view.Page(ctx))
}

func (c *ctrl) postSignIn(ctx *hime.Context) error {
	s := appctx.GetSession(ctx)
	f := s.Flash()

	email := ctx.FormValueTrimSpace("email")
	if email == "" {
		f.Add("Errors", "email required")
		return ctx.RedirectToGet()
	}

	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		f.Set("Email", email)
		f.Add("Errors", "invalid email")
		return ctx.RedirectToGet()
	}

	ok, err := repository.CanAcquireMagicLink(ctx, email)
	if err != nil {
		return err
	}
	if !ok {
		f.Add("Errors", "อีเมลของคุณได้ขอ Magic Link จากเราไปแล้ว กรุณาตรวจสอบอีเมล")
		return ctx.RedirectToGet()
	}

	f.Set("CheckEmail", "1")

	user, err := repository.GetEmailSignInUserByEmail(ctx, email)
	if err == entity.ErrNotFound {
		return ctx.RedirectTo("auth.signin.check-email")
	}
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	linkID := generateMagicLinkID()

	err = repository.StoreMagicLink(ctx, linkID, user.ID)
	if err != nil {
		return err
	}

	linkQuery := make(url.Values)
	linkQuery.Set("id", linkID)

	message := fmt.Sprintf(`สวัสดีครับคุณ %s,


ตามที่ท่านได้ขอ Magic Link เพื่อเข้าสู่ระบบสำหรับ acourse.io นั้นท่านสามารถเข้าได้ผ่าน Link ข้างล่างนี้ ภายใน 1 ชม.

%s

ทีมงาน acourse.io
	`, user.Name, c.BaseURL+ctx.Route("auth.signin.link", ctx.Param("id", linkID)))

	go c.EmailSender.Send(user.Email, "Magic Link Request", view.Markdown(message))

	return ctx.RedirectTo("auth.signin.check-email")
}

func (c *ctrl) checkEmail(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()
	if !f.Has("CheckEmail") {
		return ctx.Redirect("/")
	}
	return ctx.View("auth.check-email", view.Page(ctx))
}

func (c *ctrl) signInLink(ctx *hime.Context) error {
	linkID := ctx.FormValue("id")
	if len(linkID) == 0 {
		return ctx.RedirectTo("auth.signin")
	}

	s := appctx.GetSession(ctx)
	f := s.Flash()

	userID, err := repository.FindMagicLink(ctx, linkID)
	if err != nil {
		f.Add("Errors", "ไม่พบ Magic Link ของคุณ")
		return ctx.RedirectTo("auth.signin")
	}

	appsess.SetUserID(s, userID)
	return ctx.Redirect("/")
}

func (c *ctrl) signInPassword(ctx *hime.Context) error {
	return ctx.View("auth.signin-password", view.Page(ctx))
}

func (c *ctrl) postSignInPassword(ctx *hime.Context) error {
	s := appctx.GetSession(ctx)
	f := s.Flash()

	email := ctx.PostFormValue("email")
	if email == "" {
		f.Add("Errors", "email required")
	}
	pass := ctx.PostFormValue("password")
	if len(pass) == 0 {
		f.Add("Errors", "password required")
	}
	if f.Has("Errors") {
		f.Set("Email", email)
		return ctx.RedirectToGet()
	}

	userID, err := c.Auth.VerifyPassword(ctx, email, pass)
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	s.Regenerate()
	appsess.SetUserID(s, userID)

	// if user not found in our database, insert new user
	// this happend when database out of sync with firebase authentication
	{
		ok, err := repository.IsUserExists(ctx, userID)
		if err != nil {
			return err
		}

		email, _ = govalidator.NormalizeEmail(email)
		if !ok {
			err = repository.RegisterUser(ctx, &entity.RegisterUser{
				ID:       userID,
				Username: userID,
				Name:     userID,
				Email:    email,
			})
			if err != nil {
				return err
			}
		}
	}

	return ctx.SafeRedirect(ctx.FormValue("r"))
}
