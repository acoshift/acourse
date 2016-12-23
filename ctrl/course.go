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
	// try get by id first
	x, err := c.db.CourseGet(ctx.CourseID)
	if err != nil {
		return err
	}
	// try get by url
	if x == nil {
		x, err = c.db.CourseFind(ctx.CourseID)
		if err != nil {
			return err
		}
	}
	if x == nil {
		return ctx.NotFound()
	}

	// get owner
	owner, err := c.db.UserGet(x.Owner)
	if err != nil {
		return err
	}
	if owner == nil {
		return app.CreateError(500, "course", "can not find owner")
	}

	// get student count
	student, err := c.db.EnrollCourseCount(x.ID)
	if err != nil {
		return err
	}

	// get current user enroll
	enroll, err := c.db.EnrollFind(ctx.CurrentUserID, x.ID)
	if err != nil {
		return err
	}

	if enroll != nil {
		return ctx.OK(ToCourseView(x, ToUserTinyView(owner), student, true))
	}
	return ctx.OKPublic(ToCoursePublicView(x, ToUserTinyView(owner), student))
}

// Update runs update action
func (c *CourseController) Update(ctx *app.CourseUpdateContext) error {
	return nil
}

// List runs list action
func (c *CourseController) List(ctx *app.CourseListContext) error {
	xs, err := c.db.CourseList(store.CourseListOptionPublic(true))
	if err != nil {
		return err
	}

	res := make(app.CourseTinyCollectionView, len(xs))
	for i, x := range xs {
		u, err := c.db.UserGet(x.Owner)
		if err != nil {
			return err
		}
		if u == nil {
			return app.CreateError(500, "course", "can not find owner")
		}
		student, err := c.db.EnrollCourseCount(x.ID)
		if err != nil {
			return err
		}
		res[i] = ToCourseTinyView(x, ToUserTinyView(u), student)
	}
	return ctx.OKTiny(res)
}
