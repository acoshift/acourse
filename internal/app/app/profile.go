package app

import (
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

func signOut(ctx *hime.Context) error {
	appctx.DestroySession(ctx)
	return ctx.Redirect("/")
}

func getProfile(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)

	ownCourses, err := listOwnCourses(ctx, user.ID)
	if err != nil {
		return err
	}

	enrolledCourses, err := listEnrolledCourses(ctx, user.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Meta.Title = user.Username
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

	image, _ := ctx.FormFileHeaderNotEmpty("image")
	err := dispatcher.Dispatch(ctx, &user.UpdateProfile{
		ID:       appctx.GetUserID(ctx),
		Username: username,
		Name:     name,
		AboutMe:  aboutMe,
		Image:    image,
	})
	if app.IsUIError(err) {
		f.Add("Errors", err.Error())
		return ctx.RedirectBackToGet()
	}

	f.Clear()

	return ctx.RedirectTo("app.profile")
}
