package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/prefixhandler"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
)

type courseURLKey struct{}

func courseView(ctx *hime.Context) error {
	if ctx.Request().URL.Path != "/" {
		return notFound(ctx)
	}

	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String)
	}

	enrolled := false
	pendingEnroll := false
	if user != nil {
		enrolled, err = repository.IsEnrolled(ctx, user.ID, x.ID)
		if err != nil {
			return err
		}

		if !enrolled {
			pendingEnroll, err = repository.HasPendingPayment(ctx, user.ID, x.ID)
			if err != nil {
				return err
			}
		}
	}

	var owned bool
	if user != nil {
		owned = user.ID == x.UserID
	}

	// if user enrolled or user is owner fetch course contents
	if enrolled || owned {
		x.Contents, err = repository.GetCourseContents(ctx, x.ID)
		if err != nil {
			return err
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = repository.GetUser(ctx, x.UserID)
		if err != nil {
			return err
		}
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	page["Enrolled"] = enrolled
	page["Owned"] = owned
	page["pendingEnroll"] = pendingEnroll
	return ctx.View("course", page)
}

func courseContent(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String, "content")
	}

	enrolled, err := repository.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	x.Contents, err = repository.GetCourseContents(ctx, x.ID)
	if err != nil {
		return err
	}

	x.Owner, err = repository.GetUser(ctx, x.UserID)
	if err != nil {
		return err
	}

	var content *entity.CourseContent
	p, _ := strconv.Atoi(ctx.FormValue("p"))
	if p < 0 {
		p = 0
	}
	if p > len(x.Contents)-1 {
		p = len(x.Contents) - 1
	}
	if p >= 0 {
		content = x.Contents[p]
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["Course"] = x
	page["Content"] = content
	return ctx.View("course.content", page)
}

func editorCreate(ctx *hime.Context) error {
	return ctx.View("editor.create", newPage(ctx))
}

func postEditorCreate(ctx *hime.Context) error {
	f := appctx.GetSession(ctx).Flash()
	user := appctx.GetUser(ctx)

	var (
		title     = ctx.FormValue("Title")
		shortDesc = ctx.FormValue("ShortDesc")
		desc      = ctx.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("Image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadCourseCoverImage(ctx, image)
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
		return ctx.RedirectTo("course", id)
	}
	return ctx.RedirectTo("course", link)
}

func editorCourse(ctx *hime.Context) error {
	id := ctx.FormValue("id")
	course, err := repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.course", page)
}

func postEditorCourse(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	f := appctx.GetSession(ctx).Flash()

	var (
		title     = ctx.FormValue("Title")
		shortDesc = ctx.FormValue("ShortDesc")
		desc      = ctx.FormValue("Desc")
		imageURL  string
		start     pq.NullTime
		// assignment, _ = strconv.ParseBool(ctx.FormValue("Assignment"))
	)
	if len(title) == 0 {
		f.Add("Errors", "title required")
		return ctx.RedirectToGet()
	}

	if v := ctx.FormValue("Start"); len(v) > 0 {
		t, _ := time.Parse("2006-01-02", v)
		if !t.IsZero() {
			start.Time = t
			start.Valid = true
		}
	}

	if image, info, err := ctx.FormFileNotEmpty("Image"); err != http.ErrMissingFile {
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}

		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadCourseCoverImage(ctx, image)
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
		return ctx.RedirectTo("course", id)
	}
	return ctx.RedirectTo("course", link)
}

func editorContent(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}
	course.Contents, err = repository.GetCourseContents(ctx, id)
	if err != nil {
		return err
	}

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.content", page)
}

func postEditorContent(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	if ctx.FormValue("action") == "delete" {
		contentID := ctx.FormValue("contentId")

		err := repository.DeleteCourseContent(ctx, id, contentID)
		if err != nil {
			return err
		}
	}
	return ctx.RedirectToGet()
}

func courseEnroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("course", link)
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	return ctx.View("course.enroll", page)
}

func postCourseEnroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("course", link)
	}

	originalPrice := x.Price
	if x.Option.Discount {
		originalPrice = x.Discount
	}

	price, _ := strconv.ParseFloat(ctx.FormValue("Price"), 64)

	if price < 0 {
		f.Add("Errors", "price can not be negative")
		return ctx.RedirectToGet()
	}

	var imageURL string
	if originalPrice != 0 {
		image, info, err := ctx.FormFileNotEmpty("Image")
		if err == http.ErrMissingFile {
			f.Add("Errors", "image required")
			return ctx.RedirectToGet()
		}
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
		if !strings.Contains(info.Header.Get(header.ContentType), "image") {
			f.Add("Errors", "file is not an image")
			return ctx.RedirectToGet()
		}

		imageURL, err = uploadPaymentImage(ctx, image)
		if err != nil {
			f.Add("Errors", err.Error())
			return ctx.RedirectToGet()
		}
	}

	newPayment := false

	err = sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		if x.Price == 0 {
			err := repository.Enroll(ctx, user.ID, x.ID)
			if err != nil {
				return err
			}
		} else {
			err = repository.CreatePayment(ctx, &entity.Payment{
				CourseID:      x.ID,
				UserID:        user.ID,
				Image:         imageURL,
				Price:         price,
				OriginalPrice: originalPrice,
			})
			if err != nil {
				return err
			}
			newPayment = true
		}
		return nil
	})
	if err != nil {
		return err
	}

	if newPayment {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			sendSlackMessage(ctx, fmt.Sprintf("New payment for course %s, price %.2f", x.Title, price))
		}()
	}

	return ctx.RedirectTo("course", link)
}

func courseAssignment(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to int64 get course from id
	id := link
	_, err := uuid.Parse(link)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return notFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("course", x.URL.String, "assignment")
	}

	enrolled, err := repository.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := repository.GetAssignments(ctx, x.ID)
	if err != nil {
		return err
	}

	page := newPage(ctx)
	page["Title"] = x.Title + " | " + page["Title"].(string)
	page["Desc"] = x.ShortDesc
	page["Image"] = x.Image
	page["URL"] = baseURL + "/course/" + url.PathEscape(x.Link())
	page["Course"] = x
	page["Assignments"] = assignments
	return ctx.View("assignment", page)
}

func editorContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	course, err := repository.GetCourse(ctx, id)
	if err != nil {
		return err
	}

	page := newPage(ctx)
	page["Course"] = course
	return ctx.View("editor.content.create", page)
}

func postEditorContentCreate(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	var (
		title   = ctx.FormValue("Title")
		desc    = ctx.FormValue("Desc")
		videoID = ctx.FormValue("VideoID")
	)

	_, err := repository.RegisterCourseContent(ctx, &entity.RegisterCourseContent{
		CourseID:  id,
		Title:     title,
		LongDesc:  desc,
		VideoID:   videoID,
		VideoType: entity.Youtube,
	})
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", ctx.FormValue("id")))
}

func editorContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := repository.GetCourseContent(ctx, id)
	if err == sql.ErrNoRows {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	page := newPage(ctx)
	page["Course"] = course
	page["Content"] = content
	return ctx.View("editor.content.edit", page)
}

func postEditorContentEdit(ctx *hime.Context) error {
	// course content id
	id := ctx.FormValue("id")

	content, err := repository.GetCourseContent(ctx, id)
	if err == sql.ErrNoRows {
		return notFound(ctx)
	}
	if err != nil {
		return err
	}

	course, err := repository.GetCourse(ctx, content.CourseID)
	if err != nil {
		return err
	}

	user := appctx.GetUser(ctx)
	// user is not course owner
	if user.ID != course.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	var (
		title   = ctx.FormValue("Title")
		desc    = ctx.FormValue("Desc")
		videoID = ctx.FormValue("VideoID")
	)

	err = repository.UpdateCourseContent(ctx, course.ID, id, title, desc, videoID)
	if err != nil {
		return err
	}

	return ctx.RedirectTo("editor.content", ctx.Param("id", course.ID))
}
