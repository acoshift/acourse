package ctrl

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/gotcha"
)

// RenderController implements RenderController interface
type RenderController struct {
	db            *store.DB
	courseService acourse.CourseServiceClient
}

// NewRenderController creates controller
func NewRenderController(db *store.DB, courseService acourse.CourseServiceClient) *RenderController {
	return &RenderController{db, courseService}
}

// RenderIndex type
type RenderIndex struct {
	Title       string
	Description string
	Image       string
	URL         string
	State       map[string]interface{}
}

var cacheRender = gotcha.New()

// Index runs index action
func (c *RenderController) Index(ctx *app.RenderIndexContext) (interface{}, error) {
	var res *acourse.CoursesResponse
	refill := func() {
		xs, _ := c.courseService.ListPublicCourses(context.Background(), &acourse.ListRequest{})
		cacheRender.SetTTL("index", xs, time.Second*30)
	}
	if cache := cacheRender.Get("index"); cache != nil {
		res = cache.(*acourse.CoursesResponse)
	} else {
		// do not wait for api call
		go refill()
	}

	return &RenderIndex{
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
	response, err := c.courseService.GetCourse(context.Background(), &acourse.CourseIDRequest{CourseId: ctx.CourseID})
	if err != nil || response == nil {
		return nil, nil
	}
	course := response.Course
	r := &RenderIndex{
		Title:       course.Title,
		Description: course.ShortDescription,
		Image:       course.Photo,
		URL:         "https://acourse.io/course/",
		State: map[string]interface{}{
			"course": response,
		},
	}
	if course.Url != "" {
		r.URL += course.Url
	} else {
		r.URL += course.Id
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
