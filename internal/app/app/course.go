package app

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/prefixhandler"
	"github.com/moonrhythm/hime"
	"github.com/satori/go.uuid"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/course"
	"github.com/acoshift/acourse/internal/pkg/model/user"
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
			courseID, err = getCourseIDByURL(ctx, link)
			if err == model.ErrNotFound {
				return view.NotFound(ctx)
			}
			if err != nil {
				return err
			}
		}

		x, err := getCourse(ctx, courseID)
		if err == model.ErrNotFound {
			return view.NotFound(ctx)
		}
		if err != nil {
			return err
		}

		// if course has url, redirect to course url
		if l := x.Link(); l != link {
			return ctx.RedirectTo("app.course", l, ctx.URL.Path)
		}

		ctx = ctx.WithValue(courseKey{}, x)

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
	if ctx.URL.Path != "/" {
		return view.NotFound(ctx)
	}

	u := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	enrolled := false
	pendingEnroll := false
	if u != nil {
		enrolled := user.IsEnroll{ID: u.ID, CourseID: course.ID}
		err := dispatcher.Dispatch(ctx, &enrolled)
		if err != nil {
			return err
		}

		if !enrolled.Result {
			pendingEnroll, err = hasPendingPayment(ctx, u.ID, course.ID)
			if err != nil {
				return err
			}
		}
	}

	var owned bool
	if u != nil {
		owned = u.ID == course.Owner.ID
	}

	p := view.Page(ctx)
	p.Meta.Title = course.Title
	p.Meta.Desc = course.ShortDesc
	p.Meta.Image = course.Image
	p.Meta.URL = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p.Data["Course"] = course
	p.Data["Enrolled"] = enrolled
	p.Data["Owned"] = owned
	p.Data["PendingEnroll"] = pendingEnroll
	return ctx.View("app.course", p)
}

func (c *courseCtrl) content(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	x := c.getCourse(ctx)

	enrolled := user.IsEnroll{ID: u.ID, CourseID: x.ID}
	err := dispatcher.Dispatch(ctx, &enrolled)
	if err != nil {
		return err
	}

	if !enrolled.Result && u.ID != x.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	contents, err := getCourseContents(ctx, x.ID)
	if err != nil {
		return err
	}

	var content *course.Content
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
	p.Meta.Title = x.Title
	p.Meta.Desc = x.ShortDesc
	p.Meta.Image = x.Image
	p.Data["Course"] = x
	p.Data["Contents"] = contents
	p.Data["Content"] = content
	return ctx.View("app.course-content", p)
}

func (c *courseCtrl) enroll(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	// owner redirect to course content
	if u != nil && u.ID == course.Owner.ID {
		return ctx.RedirectTo("app.course", course.Link(), "content")
	}

	// redirect enrolled user to course content page
	enrolled := user.IsEnroll{ID: u.ID, CourseID: course.ID}
	err := dispatcher.Dispatch(ctx, &enrolled)
	if err != nil {
		return err
	}
	if enrolled.Result {
		return ctx.RedirectTo("app.course", course.Link(), "content")
	}

	// check is user has pending enroll
	pendingPayment, err := hasPendingPayment(ctx, u.ID, course.ID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", course.Link())
	}

	p := view.Page(ctx)
	p.Meta.Title = course.Title
	p.Meta.Desc = course.ShortDesc
	p.Meta.Image = course.Image
	p.Meta.URL = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p.Data["Course"] = course
	return ctx.View("app.course-enroll", p)
}

func (c *courseCtrl) postEnroll(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	x := c.getCourse(ctx)

	// owner redirect to course content
	if u != nil && u.ID == x.Owner.ID {
		return ctx.RedirectTo("app.course", x.Link(), "content")
	}

	// redirect enrolled user to course content page
	enrolled := user.IsEnroll{ID: u.ID, CourseID: x.ID}
	err := dispatcher.Dispatch(ctx, &enrolled)
	if err != nil {
		return err
	}
	if enrolled.Result {
		return ctx.RedirectTo("app.course", x.Link(), "content")
	}

	// check is user has pending enroll
	pendingPayment, err := hasPendingPayment(ctx, u.ID, x.ID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", x.Link())
	}

	f := appctx.GetFlash(ctx)

	price, _ := strconv.ParseFloat(ctx.FormValue("price"), 64)
	image, _ := ctx.FormFileHeaderNotEmpty("image")

	err = dispatcher.Dispatch(ctx, &user.Enroll{
		ID:           u.ID,
		CourseID:     x.ID,
		Price:        price,
		PaymentImage: image,
	})
	if app.IsUIError(err) {
		f.Add("Errors", "image required")
		return ctx.RedirectToGet()
	}
	if err != nil {
		return err
	}

	return ctx.RedirectTo("app.course", x.Link())
}

func (c *courseCtrl) assignment(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	course := c.getCourse(ctx)

	enrolled := user.IsEnroll{ID: u.ID, CourseID: course.ID}
	err := dispatcher.Dispatch(ctx, &enrolled)
	if err != nil {
		return err
	}

	if !enrolled.Result && u.ID != course.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := findAssignmentsByCourseID(ctx, course.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Meta.Title = course.Title
	p.Meta.Desc = course.ShortDesc
	p.Meta.Image = course.Image
	p.Meta.URL = c.BaseURL + ctx.Route("app.course", url.PathEscape(course.Link()))
	p.Data["Course"] = course
	p.Data["Assignments"] = assignments
	return ctx.View("app.course-assignment", p)
}
