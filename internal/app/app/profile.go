package app

import (
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/image"
	"github.com/acoshift/acourse/internal/pkg/me"
)

func signOut(ctx *hime.Context) error {
	appctx.DestroySession(ctx)
	return ctx.Redirect("/")
}

func getProfile(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)

	ownCourses, err := me.GetOwnCourses(ctx, u.ID)
	if err != nil {
		return err
	}

	enrolledCourses, err := me.GetEnrolledCourses(ctx, u.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Meta.Title = u.Username
	p.Data["Navbar"] = "profile"
	p.Data["OwnCourses"] = ownCourses
	p.Data["EnrolledCourses"] = enrolledCourses
	return ctx.View("app.profile", p)
}

func getProfileEdit(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	f := appctx.GetFlash(ctx)
	if !f.Has("Username") {
		f.Set("Username", user.Username)
	}
	if !f.Has("Name") {
		f.Set("Name", user.Name)
	}
	if !f.Has("AboutMe") {
		f.Set("AboutMe", user.AboutMe)
	}

	p := view.Page(ctx)
	p.Meta.Title = user.Username
	return ctx.View("app.profile-edit", p)
}

func postProfileEdit(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	var (
		username = ctx.FormValue("username")
		name     = ctx.FormValue("name")
		aboutMe  = ctx.FormValue("aboutMe")
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

	img, _ := ctx.FormFileHeaderNotEmpty("image")
	err := me.UpdateProfile(ctx, &me.UpdateProfileArgs{
		Username: username,
		Name:     name,
		AboutMe:  aboutMe,
		Image:    img,
	})
	if err == image.ErrInvalidType {
		f.Add("Errors", "invalid image")
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	f.Clear()

	return ctx.RedirectTo("app.profile")
}
