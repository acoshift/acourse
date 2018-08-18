package app

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/prefixhandler"
	"github.com/satori/go.uuid"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

type courseURLKey struct{}

func (c *ctrl) courseView(ctx *hime.Context) error {
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
		id, err = c.Repository.GetCourseIDByURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := c.Repository.GetCourse(ctx, id)
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
		enrolled, err = c.Repository.IsEnrolled(ctx, user.ID, x.ID)
		if err != nil {
			return err
		}

		if !enrolled {
			pendingEnroll, err = c.Repository.HasPendingPayment(ctx, user.ID, x.ID)
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
		x.Contents, err = c.Repository.GetCourseContents(ctx, x.ID)
		if err != nil {
			return err
		}
	}

	if owned {
		x.Owner = user
	} else {
		x.Owner, err = c.Repository.GetUser(ctx, x.UserID)
		if err != nil {
			return err
		}
	}

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	p["Enrolled"] = enrolled
	p["Owned"] = owned
	p["PendingEnroll"] = pendingEnroll
	return ctx.View("app.course", p)
}

func (c *ctrl) courseContent(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to uuid get course from id
	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		// link can not parse to uuid get course id from url
		id, err = c.Repository.GetCourseIDByURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := c.Repository.GetCourse(ctx, id)
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

	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	x.Contents, err = c.Repository.GetCourseContents(ctx, x.ID)
	if err != nil {
		return err
	}

	x.Owner, err = c.Repository.GetUser(ctx, x.UserID)
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

func (c *ctrl) courseEnroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		id, err = c.Repository.GetCourseIDByURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	x, err := c.Repository.GetCourse(ctx, id)
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
	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, id)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.course", link)
	}

	// check is user has pending enroll
	pendingPayment, err := c.Repository.HasPendingPayment(ctx, user.ID, id)
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
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	return ctx.View("app.course-enroll", p)
}

func (c *ctrl) postCourseEnroll(ctx *hime.Context) error {
	f := appctx.GetFlash(ctx)

	link := prefixhandler.Get(ctx, courseURLKey{})

	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		id, err = c.Repository.GetCourseIDByURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}

	_, err = c.Repository.GetCourse(ctx, id)
	if err == entity.ErrNotFound {
		return share.NotFound(ctx)
	}
	if err != nil {
		return err
	}

	price, _ := strconv.ParseFloat(ctx.FormValue("price"), 64)
	image, _ := ctx.FormFileHeaderNotEmpty("image")

	err = c.Service.EnrollCourse(ctx, id, price, image)
	if service.IsUIError(err) {
		f.Add("Errors", "image required")
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectTo("app.course", link)
}

func (c *ctrl) courseAssignment(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	link := prefixhandler.Get(ctx, courseURLKey{})

	// if id can parse to int64 get course from id
	id := link
	_, err := uuid.FromString(link)
	if err != nil {
		// link can not parse to int64 get course id from url
		id, err = c.Repository.GetCourseIDByURL(ctx, link)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}
	}
	x, err := c.Repository.GetCourse(ctx, id)
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

	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, x.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != x.UserID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := c.Repository.GetAssignments(ctx, x.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Title"] = x.Title
	p["Desc"] = x.ShortDesc
	p["Image"] = x.Image
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(x.Link()))
	p["Course"] = x
	p["Assignments"] = assignments
	return ctx.View("app.course-assignment", p)
}
