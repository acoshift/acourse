package ctrl

import (
	"time"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/acourse/pkg/view"
)

// RenderController implements RenderController interface
type RenderController struct {
	db         *store.DB
	courseCtrl app.CourseController
}

// NewRenderController creates controller
func NewRenderController(db *store.DB, courseCtrl app.CourseController) *RenderController {
	return &RenderController{db, courseCtrl}
}

var cacheRender = store.NewCache(time.Second * 15)

// Index runs index action
func (c *RenderController) Index(ctx *app.RenderIndexContext) (interface{}, error) {
	var res view.CourseTinyCollection
	if cache := cacheRender.Get("index"); cache != nil {
		res = cache.(view.CourseTinyCollection)
	} else {
		// do not wait for api call
		go func() {
			xs, _ := c.courseCtrl.List(&app.CourseListContext{})
			cacheRender.Set("index", xs)
		}()
	}

	return &view.RenderIndex{
		Title:       "Acourse",
		Description: "Online courses for everyone",
		Image:       "https://acourse.io/static/acourse-og.jpg",
		URL:         "https://acourse.io",
		State: map[string]interface{}{
			"courses": res,
		},
	}, nil
}

// Course runs course action
func (c *RenderController) Course(ctx *app.RenderCourseContext) (interface{}, error) {
	courseInf, err := c.courseCtrl.Show(&app.CourseShowContext{CourseID: ctx.CourseID})
	if err != nil || courseInf == nil {
		return nil, nil
	}
	course, ok := courseInf.(*view.CoursePublic)
	if !ok {
		return nil, nil
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

	return r, nil
}
