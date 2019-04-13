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
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/me"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

type (
	courseIDKey struct{}
	courseKey   struct{}
)

func newCourseHandler() http.Handler {
	c := courseCtrl{}

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
			courseID, err = course.GetIDByURL(ctx, link)
			if err == course.ErrNotFound {
				return view.NotFound(ctx)
			}
			if err != nil {
				return err
			}
		}

		x, err := course.Get(ctx, courseID)
		if err == course.ErrNotFound {
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

type courseCtrl struct{}

func (ctrl *courseCtrl) getCourse(ctx context.Context) *course.Course {
	return ctx.Value(courseKey{}).(*course.Course)
}

func (ctrl *courseCtrl) view(ctx *hime.Context) error {
	if ctx.URL.Path != "/" {
		return view.NotFound(ctx)
	}

	u := appctx.GetUser(ctx)
	c := ctrl.getCourse(ctx)

	enrolled := false
	pendingEnroll := false
	var err error
	if u != nil {
		enrolled, err = course.IsEnroll(ctx, u.ID, c.ID)
		if err != nil {
			return err
		}

		if !enrolled {
			pendingEnroll, err = payment.HasPending(ctx, u.ID, c.ID)
			if err != nil {
				return err
			}
		}
	}

	var owned bool
	if u != nil {
		owned = u.ID == c.Owner.ID
	}

	p := view.Page(ctx)
	p.Meta.Title = c.Title
	p.Meta.Desc = c.ShortDesc
	p.Meta.Image = c.Image
	p.Meta.URL = ctx.Global("baseURL").(string) + ctx.Route("app.course", url.PathEscape(c.Link()))
	p.Data["Course"] = c
	p.Data["Enrolled"] = enrolled
	p.Data["Owned"] = owned
	p.Data["PendingEnroll"] = pendingEnroll
	return ctx.View("app.course", p)
}

func (ctrl *courseCtrl) content(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	x := ctrl.getCourse(ctx)

	enrolled, err := course.IsEnroll(ctx, u.ID, x.ID)
	if err != nil {
		return err
	}

	if !enrolled && u.ID != x.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	contents, err := course.GetContents(ctx, x.ID)
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

func (ctrl *courseCtrl) enroll(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	c := ctrl.getCourse(ctx)

	// owner redirect to c content
	if u != nil && u.ID == c.Owner.ID {
		return ctx.RedirectTo("app.c", c.Link(), "content")
	}

	// redirect enrolled user to c content page
	enrolled, err := course.IsEnroll(ctx, u.ID, c.ID)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.c", c.Link(), "content")
	}

	// check is user has pending enroll
	pendingPayment, err := payment.HasPending(ctx, u.ID, c.ID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.c", c.Link())
	}

	p := view.Page(ctx)
	p.Meta.Title = c.Title
	p.Meta.Desc = c.ShortDesc
	p.Meta.Image = c.Image
	p.Meta.URL = ctx.Global("baseURL").(string) + ctx.Route("app.c", url.PathEscape(c.Link()))
	p.Data["Course"] = c
	return ctx.View("app.c-enroll", p)
}

func (ctrl *courseCtrl) postEnroll(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	x := ctrl.getCourse(ctx)

	// owner redirect to course content
	if u != nil && u.ID == x.Owner.ID {
		return ctx.RedirectTo("app.course", x.Link(), "content")
	}

	// redirect enrolled user to course content page
	enrolled, err := course.IsEnroll(ctx, u.ID, x.ID)
	if err != nil {
		return err
	}
	if enrolled {
		return ctx.RedirectTo("app.course", x.Link(), "content")
	}

	// check is user has pending enroll
	pendingPayment, err := payment.HasPending(ctx, u.ID, x.ID)
	if err != nil {
		return err
	}
	if pendingPayment {
		return ctx.RedirectTo("app.course", x.Link())
	}

	f := appctx.GetFlash(ctx)

	price, _ := strconv.ParseFloat(ctx.FormValue("price"), 64)
	image, _ := ctx.FormFileHeaderNotEmpty("image")

	if price < 0 {
		f.Add("Errors", "จำนวนเงินติดลบไม่ได้")
		return ctx.RedirectToGet()
	}

	err = me.Enroll(ctx, x.ID, price, image)
	if err == me.ErrImageRequired {
		f.Add("Errors", "กรุณาอัพโหลดรูปภาพ")
		return ctx.RedirectToGet()
	}
	if err != nil {
		f.Add("Errors", "image required")
		return ctx.RedirectToGet()
	}

	return ctx.RedirectTo("app.course", x.Link())
}

func (ctrl *courseCtrl) assignment(ctx *hime.Context) error {
	u := appctx.GetUser(ctx)
	c := ctrl.getCourse(ctx)

	enrolled, err := course.IsEnroll(ctx, u.ID, c.ID)
	if err != nil {
		return err
	}

	if !enrolled && u.ID != c.Owner.ID {
		return ctx.Status(http.StatusForbidden).StatusText()
	}

	assignments, err := course.GetAssignments(ctx, c.ID)
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Meta.Title = c.Title
	p.Meta.Desc = c.ShortDesc
	p.Meta.Image = c.Image
	p.Meta.URL = ctx.Global("baseURL").(string) + ctx.Route("app.c", url.PathEscape(c.Link()))
	p.Data["Course"] = c
	p.Data["Assignments"] = assignments
	return ctx.View("app.c-assignment", p)
}
