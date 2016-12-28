package ctrl

import (
	"acourse/app"
	"acourse/store"
	"acourse/view"
	"time"
)

// RenderController implements RenderController interface
type RenderController struct {
	db *store.DB
}

// NewRenderController creates controller
func NewRenderController(db *store.DB) *RenderController {
	return &RenderController{db}
}

var cacheRender = store.NewCache(time.Second * 15)

// Index runs index action
func (c *RenderController) Index(ctx *app.RenderIndexContext) error {
	var res view.CourseTinyCollection
	if cache := cacheRender.Get("index"); cache != nil {
		res = cache.(view.CourseTinyCollection)
	} else {
		// do not wait for api call
		go func() {
			xs, _ := c.db.CourseList(store.CourseListOptionPublic(true))
			rs := make(view.CourseTinyCollection, len(xs))
			for i, x := range xs {
				u, _ := c.db.UserMustGet(x.Owner)
				student, _ := c.db.EnrollCourseCount(x.ID)
				rs[i] = ToCourseTinyView(x, ToUserTinyView(u), student)
			}
			cacheRender.Set("index", rs)
		}()
	}

	return ctx.OK(&view.RenderIndex{
		Title:       "Acourse",
		Description: "Online courses for everyone",
		Image:       "https://acourse.io/static/acourse-og.jpg",
		URL:         "https://acourse.io",
		State: map[string]interface{}{
			"courses": res,
		},
	})
}

// Course runs course action
func (c *RenderController) Course(ctx *app.RenderCourseContext) error {
	course, err := c.db.CourseFind(ctx.CourseID)
	if course == nil {
		course, err = c.db.CourseGet(ctx.CourseID)
	}
	if err != nil || course == nil {
		return ctx.NotFound()
	}
	r := &view.RenderIndex{
		Title:       course.Title,
		Description: course.ShortDescription,
		Image:       course.Photo,
		URL:         "https://acourse.io/course/",
	}
	if course.URL != "" {
		r.URL += course.URL
	} else {
		r.URL += course.ID
	}

	if r.Title == "" {
		r.Title = "Acourse"
	} else {
		r.Title += " | Acourse"
	}
	if r.Description == "" {
		r.Description = "Online courses for everyone"
	}
	if r.Image == "" {
		r.Image = "https://acourse.io/static/acourse-og.jpg"
	}

	return ctx.OK(r)
}
