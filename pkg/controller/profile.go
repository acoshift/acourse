package controller

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/header"
	"github.com/acoshift/session"
	"github.com/asaskevich/govalidator"
)

func (c *ctrl) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(r.Context())

	ownCourses, err := c.repo.ListOwnCourses(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enrolledCourses, err := c.repo.ListEnrolledCourses(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.Profile(w, r, ownCourses, enrolledCourses)
}

func (c *ctrl) ProfileEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postProfileEdit(w, r)
		return
	}
	ctx := r.Context()
	user := app.GetUser(ctx)
	f := session.Get(ctx, sessName).Flash()
	if !f.Has("Username") {
		f.Set("Username", user.Username)
	}
	if !f.Has("Name") {
		f.Set("Name", user.Name)
	}
	if !f.Has("AboutMe") {
		f.Set("AboutMe", user.AboutMe)
	}
	c.view.ProfileEdit(w, r)
}

func (c *ctrl) postProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := app.GetUser(ctx)
	f := session.Get(ctx, sessName).Flash()

	image, info, err := r.FormFile("Image")
	var imageURL string
	if err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			back(w, r)
			return
		}

		imageURL, err = UploadProfileImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			back(w, r)
			return
		}
	}

	var (
		username = r.FormValue("Username")
		name     = r.FormValue("Name")
		aboutMe  = r.FormValue("AboutMe")
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
		back(w, r)
		return
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if len(imageURL) > 0 {
		_, err = tx.Exec(`
			update users
			set image = $2
			where id = $1
		`, user.ID, imageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, err = tx.Exec(`
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, user.ID, username, name, aboutMe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		f.Add("Errors", err.Error())
		back(w, r)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
