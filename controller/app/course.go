package app

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/satori/go.uuid"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/service"
	"github.com/acoshift/acourse/view"
)

type (
	courseIDKey struct{}
	courseKey   struct{}
)

func newCourseHandler(appCtrl *ctrl) http.Handler {
	c := courseCtrl{
		ctrl: appCtrl,
	}

	mux := http.NewServeMux()
	mux.Handle("/", methodmux.Get(
		hime.Handler(c.view),
	))
	mux.Handle("/content", mustSignedIn(methodmux.Get(
		hime.Handler(c.content),
	)))
	mux.Handle("/enroll", mustSignedIn(methodmux.GetPost(
		hime.Handler(c.enroll),
		hime.Handler(c.postEnroll),
	)))
	mux.Handle("/assignment", mustSignedIn(methodmux.Get(
		hime.Handler(c.assignment),
	)))

	return hime.Handler(func(ctx *hime.Context) error {
		link := prefixhandler.Get(ctx, courseIDKey{})

		courseID := link
		_, err := uuid.FromString(link)
		if err != nil {
			// link can not parse to uuid get course id from url
			courseID, err = c.Repository.GetCourseIDByURL(ctx, link)
			if err == entity.ErrNotFound {
				return share.NotFound(ctx)
			}
			if err != nil {
				return err
			}
		}

		x, err := c.Repository.GetCourse(ctx, courseID)
		if err == entity.ErrNotFound {
			return share.NotFound(ctx)
		}
		if err != nil {
			return err
		}

		// if course has url, redirect to course url
		if l := x.Link(); l != link {
			return ctx.RedirectTo("app.course", l)
		}

		ctx.WithValue(courseKey{}, x)

		return ctx.Handle(mux)
	})
}

type courseCtrl struct {
	*ctrl
}

func (c *courseCtrl) getCourse(ctx context.Context) *Course {
	return ctx.Value(courseKey{}).(*Course)
}

func (c *courseCtrl) view(ctx *hime.Context) error {
	if ctx.Request().URL.Path != "/" {
		return share.NotFound(ctx)
	}

	user := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	var err error
	enrolled := false
	pendingEnroll := false
	if user != nil {
		enrolled, err = c.Repository.IsEnrolled(ctx, user.ID, course.ID)
		if err != nil {
			return err
		}

		if !enrolled {
			pendingEnroll, err = c.Repository.HasPendingPayment(ctx, user.ID, course.ID)
			if err != nil {
				return err
			}
		}
	}

	var owned bool
	if user != nil {
		owned = user.ID == course.Owner.ID
	}

	p := view.Page(ctx)
	p["Title"] = course.Title
	p["Desc"] = course.ShortDesc
	p["Image"] = course.Image
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p["Course"] = course
	p["Enrolled"] = enrolled
	p["Owned"] = owned
	p["PendingEnroll"] = pendingEnroll
	return ctx.View("app.course", p)
}

func (c *courseCtrl) content(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, course.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != course.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	contents, err := c.Repository.GetCourseContents(ctx, course.ID)
	if err != nil {
		return err
	}

	var content *entity.CourseContent
	pg, _ := strconv.Atoi(ctx.FormValue("p"))
	if pg < 0 {
		pg = 0
	}
	if pg > len(contents)-1 {
		pg = len(contents) - 1
	}
	if pg >= 0 {
		content = contents[pg]
	}

	p := view.Page(ctx)
	p["Title"] = course.Title
	p["Desc"] = course.ShortDesc
	p["Image"] = course.Image
	p["Course"] = course
	p["Contents"] = contents
	p["Content"] = content
	return ctx.View("app.course-content", p)
}

func (c *courseCtrl) enroll(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	// redirect enrolled user back to course page
	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, course.ID)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.course", course.Link())
	}

	// check is user has pending enroll
	pendingPayment, err := c.Repository.HasPendingPayment(ctx, user.ID, course.ID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", course.Link())
	}

	p := view.Page(ctx)
	p["Title"] = course.Title
	p["Desc"] = course.ShortDesc
	p["Image"] = course.Image
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p["Course"] = course
	return ctx.View("app.course-enroll", p)
}

func (c *courseCtrl) postEnroll(ctx *hime.Context) error {
	course := c.getCourse(ctx)

	f := appctx.GetFlash(ctx)

	price, _ := strconv.ParseFloat(ctx.FormValue("price"), 64)
	image, _ := ctx.FormFileHeaderNotEmpty("image")

	err := c.Service.EnrollCourse(ctx, course.ID, price, image)
	if service.IsUIError(err) {
		f.Add("Errors", "image required")
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectTo("app.course", course.Link())
}

func (c *courseCtrl) assignment(ctx *hime.Context) error {
	user := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	enrolled, err := c.Repository.IsEnrolled(ctx, user.ID, course.ID)
	if err != nil {
		return err
	}

	if !enrolled && user.ID != course.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := c.Repository.FindAssignmentsByCourseID(ctx, course.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Title"] = course.Title
	p["Desc"] = course.ShortDesc
	p["Image"] = course.Image
	p["URL"] = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p["Course"] = course
	p["Assignments"] = assignments
	return ctx.View("app.course-assignment", p)
}
