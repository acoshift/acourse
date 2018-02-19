package app

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/pgsql"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/repository"
)

func profile(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)

	ownCourses, err := repository.ListOwnCourses(db, user.ID)
	must(err)

	enrolledCourses, err := repository.ListEnrolledCourses(db, user.ID)
	must(err)

	page := newPage(ctx)
	page["Title"] = user.Username + " | " + page["Title"].(string)
	page["OwnCourses"] = ownCourses
	page["EnrolledCourses"] = enrolledCourses
	return ctx.View("profile", page)
}

func profileEdit(ctx hime.Context) hime.Result {
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

func postProfileEdit(ctx hime.Context) hime.Result {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	var imageURL string
	if image, info, err := ctx.FormFile("Image"); err != http.ErrMissingFile && info.Size > 0 {
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

	err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		if len(imageURL) > 0 {
			_, err := tx.Exec(`
				update users
				set image = $2
				where id = $1
			`, user.ID, imageURL)
			if err != nil {
				return err
			}
		}
		_, err := tx.Exec(`
			update users
			set
				username = $2,
				name = $3,
				about_me = $4,
				updated_at = now()
			where id = $1
		`, user.ID, username, name, aboutMe)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		f.Add("Errors", err.Error())
		return ctx.RedirectToGet()
	}

	return ctx.RedirectTo("profile")
}
