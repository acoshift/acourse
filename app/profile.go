package app

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/acoshift/header"
	"github.com/acoshift/pgsql"
	"github.com/asaskevich/govalidator"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func profile(w http.ResponseWriter, r *http.Request) {
	user := appctx.GetUser(r.Context())

	ownCourses, err := repository.ListOwnCourses(db, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enrolledCourses, err := repository.ListEnrolledCourses(db, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.Profile(w, r, ownCourses, enrolledCourses)
}

func profileEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postProfileEdit(w, r)
		return
	}
	ctx := r.Context()
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
	view.ProfileEdit(w, r)
}

func postProfileEdit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	var imageURL string
	if image, info, err := r.FormFile("Image"); err != http.ErrMissingFile && info.Size > 0 {
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

		imageURL, err = uploadProfileImage(ctx, image)
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
		back(w, r)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
