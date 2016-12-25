package ctrl

import (
	"acourse/app"
	"acourse/store"
)

// RenderController implements RenderController interface
type RenderController struct {
	db *store.DB
}

// NewRenderController creates controller
func NewRenderController(db *store.DB) *RenderController {
	return &RenderController{db}
}

// Index runs index action
func (c *RenderController) Index(ctx *app.RenderIndexContext) error {
	return ctx.OK(&app.RenderIndexView{
		Title:       "Acourse",
		Description: "Online courses for everyone",
		Image:       "https://acourse.io/static/acourse-og.jpg",
		URL:         "https://acourse.io",
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
	r := &app.RenderIndexView{
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
