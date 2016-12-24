package ctrl

import (
	"acourse/app"
	"acourse/store"
)

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

	if enroll != nil || ctx.CurrentUserID == x.Owner {
		return ctx.OK(ToCourseView(x, ToUserTinyView(owner), student, enroll != nil, ctx.CurrentUserID == x.Owner))
	}
	return ctx.OKPublic(ToCoursePublicView(x, ToUserTinyView(owner), student))
}

// Update runs update action
func (c *CourseController) Update(ctx *app.CourseUpdateContext) error {
	return nil
}

// List runs list action
func (c *CourseController) List(ctx *app.CourseListContext) error {
	var xs []*store.Course
	var err error

	// query with owner
	if ctx.Owner != "" {
		if ctx.Owner == ctx.CurrentUserID {
			xs, err = c.db.CourseList(store.CourseListOptionOwner(ctx.Owner))
		} else {
			xs, err = c.db.CourseList(store.CourseListOptionOwner(ctx.Owner), store.CourseListOptionPublic(true))
		}
	} else if ctx.Student != "" {
		if ctx.Student == ctx.CurrentUserID {
			var enrolls []*store.Enroll
			enrolls, err = c.db.EnrollListByUserID(ctx.Student)
			if err != nil {
				return err
			}
			ids := make([]string, len(enrolls))
			for i, e := range enrolls {
				ids[i] = e.CourseID
			}
			xs, err = c.db.CourseGetAllByIDs(ids)
		} else {
			var enrolls []*store.Enroll
			enrolls, err = c.db.EnrollListByUserID(ctx.Student)
			if err != nil {
				return err
			}
			ids := make([]string, len(enrolls))
			for i, e := range enrolls {
				ids[i] = e.CourseID
			}
			var ts []*store.Course
			ts, err = c.db.CourseGetAllByIDs(ids)
			if err != nil {
				return err
			}
			for _, t := range ts {
				if t.Options.Public {
					xs = append(xs, t)
				}
			}
		}
	} else {
		xs, err = c.db.CourseList(store.CourseListOptionPublic(true))
	}

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
