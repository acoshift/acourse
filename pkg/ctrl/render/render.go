package render

import (
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/gotcha"
)

// New creates new render controller
func New(courseService acourse.CourseServiceClient) app.RenderController {
	return &renderController{courseService}
}

// renderController implements RenderController interface
type renderController struct {
	courseService acourse.CourseServiceClient
}

type index struct {
	Title       string
	Description string
	Image       string
	URL         string
	State       map[string]interface{}
}

var cacheRender = gotcha.New()

// Index runs index action
func (c *renderController) Index(ctx *app.RenderIndexContext) (interface{}, error) {
	var res *acourse.CoursesResponse
	if cache := cacheRender.Get("index"); cache != nil {
		res = cache.(*acourse.CoursesResponse)
	} else {
		res, _ = c.courseService.ListPublicCourses(ctx, &acourse.ListRequest{})
		cacheRender.SetTTL("index", res, time.Second*30)
	}

	return &index{
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
func (c *renderController) Course(ctx *app.RenderCourseContext) (interface{}, error) {
	response, err := c.courseService.GetCourse(ctx, &acourse.CourseIDRequest{CourseId: ctx.CourseID})
	if err != nil || response == nil {
		return nil, nil
	}
	course := response.Course
	r := &index{
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
