package app

import (
	"context"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

func profile(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)

	ownCourses, err := repository.ListOwnCourses(ctx, user.ID)
	if err != nil {
		return err
	}

	enrolledCourses, err := repository.ListEnrolledCourses(ctx, user.ID)
	if err != nil {
		return err
	}

	page := newPage(ctx)
	page["Title"] = user.Username + " | " + page["Title"].(string)
	page["OwnCourses"] = ownCourses
	page["EnrolledCourses"] = enrolledCourses
	return ctx.View("profile", page)
}

func profileEdit(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()
	if !f.Has("Username") {
		f.Set("Username", user.Username)
	}
	if !f.Has("Name") {
		f.Set("Name", user.Name)
	}
	if !f.Has("AboutMe") {
		f.Set("AboutMe", user.AboutMe)
	}

	page := newPage(ctx)
	page["Title"] = user.Username + " | " + page["Title"].(string)
	return ctx.View("profile.edit", page)
}

func postProfileEdit(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	var imageURL string
	if image, info, err := ctx.FormFileNotEmpty("Image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	var (
		username = ctx.FormValue("Username")
		name     = ctx.FormValue("Name")
		aboutMe  = ctx.FormValue("AboutMe")
	)
	f.Set("Username", username)
	f.Set("Name", name)
	f.Set("AboutMe", aboutMe)

	if !govalidator.IsAlphanumeric(username) {
		f.Add("Errors", "username allow only a-z, A-Z, and 0-9")
	}
	if n := utf8.RuneCountInString(username); n < 4 || n > 32 {
		f.Add("Errors", "username must have 4 - 32 characters")
	}
	if n := utf8.RuneCountInString(name); n < 4 || n > 40 {
		f.Add("Errors", "name must have 4 - 40 characters")
	}
	if n := utf8.RuneCountInString(aboutMe); n > 256 {
		f.Add("Errors", "about me must have lower than 256 characters")
	}
	if f.Has("Errors") {
		return ctx.RedirectToGet()
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if len(imageURL) > 0 {
			err := repository.SetUserImage(ctx, user.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return repository.UpdateUser(ctx, &entity.UpdateUser{
			ID:       user.ID,
			Username: username,
			Name:     name,
			AboutMe:  aboutMe,
		})
	})
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	return ctx.RedirectTo("profile")
}
