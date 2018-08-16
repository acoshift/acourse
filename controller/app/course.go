package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/acoshift/header"
	"github.com/acoshift/hime"
	"github.com/acoshift/prefixhandler"
	"github.com/satori/go.uuid"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

type courseURLKey struct{}

func courseView(ctx *hime.Context) error {
	if ctx.Request().URL.Path != "/" {
		return share.NotFound(ctx)
	}

	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("app.course", x.URL.String)
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

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["URL"] = baseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	p["Enrolled"] = enrolled
	p["Owned"] = owned
	p["PendingEnroll"] = pendingEnroll
	return ctx.View("app.course", p)
}

func courseContent(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("app.course", x.URL.String, "content")
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
	pg, _ := strconv.Atoi(ctx.FormValue("p"))
	if pg < 0 {
		pg = 0
	}
	if pg > len(x.Contents)-1 {
		pg = len(x.Contents) - 1
	}
	if pg >= 0 {
		content = x.Contents[pg]
	}

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["Course"] = x
	p["Content"] = content
	return ctx.View("app.course-content", p)
}

func courseEnroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("app.course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", link)
	}

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["URL"] = baseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	return ctx.View("app.course-enroll", p)
}

func postCourseEnroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	f := appctx.GetSession(ctx).Flash()

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	// if user is course owner redirect back to course page
	if user.ID == x.UserID {
		return ctx.RedirectTo("app.course", link)
	}

	// redirect enrolled user back to course page
	enrolled, err := repository.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := repository.HasPendingPayment(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", link)
	}

	originalPrice := x.Price
	if x.Option.Discount {
		originalPrice = x.Discount
	}

	price, _ := strconv.ParseFloat(ctx.FormValue("price"), 64)

	if price < 0 {
		f.Add("Errors", "price can not be negative")
		return ctx.RedirectToGet()
	}

	var imageURL string
	if originalPrice != 0 {
		image, info, err := ctx.FormFileNotEmpty("image")
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
		go adminNotifier.Notify(fmt.Sprintf("New payment for course %s, price %.2f", x.Title, price))
	}

	return ctx.RedirectTo("app.course", link)
}

func courseAssignment(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to int64 get course from id
	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = repository.GetCourseIDFromURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	// if course has url, redirect to course url
	if x.URL.Valid && x.URL.String != link {
		return ctx.RedirectTo("app.course", x.URL.String, "assignment")
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

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["URL"] = baseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	p["Assignments"] = assignments
	return ctx.View("app.course-assignment", p)
}
