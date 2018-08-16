package editor

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) courseCreate(ctx *hime.Context) error {
	return ctx.View("editor.course-create", view.Page(ctx))
}

func (c *ctrl) postCourseCreate(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()
	user := appctx.GetUser(ctx)

	var (
		title     = ctx.FormValue("title")
		shortDesc = ctx.FormValue("shortDesc")
		desc      = ctx.FormValue("desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = c.uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	var id string
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		id, err = repository.RegisterCourse(ctx, &entity.RegisterCourse{
			UserID:    user.ID,
			Title:     title,
			ShortDesc: shortDesc,
			LongDesc:  desc,
			Image:     imageURL,
			Start:     start,
		})
		if err != nil {
			return err
		}

		return repository.SetCourseOption(ctx, id, &entity.CourseOption{})
	})
	if err != nil {
		return err
	}

	link, _ := repository.GetCourseURL(ctx, id)
	if link == "" {
		return ctx.RedirectTo("app.course", id)
	}
	return ctx.RedirectTo("app.course", link)
}

func (c *ctrl) courseEdit(ctx *hime.Context) error {
	id := ctx.FormValue("id")
	course, err := repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Course"] = course
	return ctx.View("editor.course-edit", p)
}

func (c *ctrl) postCourseEdit(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	f := appctx.GetSession(ctx).Flash()

	var (
		title     = ctx.FormValue("title")
		shortDesc = ctx.FormValue("shortDesc")
		desc      = ctx.FormValue("desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = c.uploadCourseCoverImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := repository.UpdateCourse(ctx, &entity.UpdateCourse{
			ID:        id,
			Title:     title,
			ShortDesc: shortDesc,
			LongDesc:  desc,
			Start:     start,
		})
		if err != nil {
			return err
		}

		if len(imageURL) > 0 {
			err = repository.SetCourseImage(ctx, id, imageURL)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	link, _ := repository.GetCourseURL(ctx, id)
	if link == "" {
		return ctx.RedirectTo("app.course", id)
	}
	return ctx.RedirectTo("app.course", link)
}
