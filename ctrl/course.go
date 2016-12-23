package ctrl

import "acourse/store"
import "acourse/app"

// CourseController implements CourseController interface
type CourseController struct {
	db *store.DB
}

// NewCourseController creates new controller
func NewCourseController(db *store.DB) *CourseController {
	return &CourseController{db: db}
}

// Show runs show action
func (c *CourseController) Show(ctx *app.CourseShowContext) error {
	return nil
}

// Update runs update action
func (c *CourseController) Update(ctx *app.CourseUpdateContext) error {
	return nil
}
