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
	role, err := c.db.RoleFindByUserID(ctx.CurrentUserID)
	if err != nil {
		return err
	}
	course, err := c.db.CourseGet(ctx.CourseID)
	if err != nil {
		return err
	}
	if course == nil {
		return ctx.NotFound()
	}
	if course.Owner != ctx.CurrentUserID || !role.Admin {
		return ctx.Forbidden()
	}

	// merge course with payload
	course.Title = ctx.Payload.Title
	course.ShortDescription = ctx.Payload.ShortDescription
	course.Description = ctx.Payload.Description
	course.Photo = ctx.Payload.Photo
	course.Start = ctx.Payload.Start
	course.Video = ctx.Payload.Video
	course.Type = store.CourseType(ctx.Payload.Type)
	course.Contents = ToCourseContents(ctx.Payload.Contents)
	course.Options.Enroll = ctx.Payload.Enroll
	course.Options.Public = ctx.Payload.Public
	course.Options.Attend = ctx.Payload.Attend
	course.Options.Assignment = ctx.Payload.Assignment
	course.Options.Purchase = ctx.Payload.Purchase

	err = c.db.CourseSave(course)
	if err != nil {
		return err
	}

	return ctx.NoContent()
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

// Enroll runs enroll action
func (c *CourseController) Enroll(ctx *app.CourseEnrollContext) error {
	course, err := c.db.CourseGet(ctx.CourseID)
	if err != nil {
		return err
	}
	if course == nil {
		return ctx.NotFound()
	}

	// owner can not enroll
	if course.Owner == ctx.CurrentUserID {
		return ctx.Forbidden()
	}

	// check is user already enrolled
	enroll, err := c.db.EnrollFind(ctx.CurrentUserID, ctx.CourseID)
	if err != nil {
		return err
	}
	if enroll != nil {
		// user already enroll
		return ctx.Forbidden()
	}

	// check is user already send waiting payment
	payment, err := c.db.PaymentFind(ctx.CurrentUserID, ctx.CourseID, store.PaymentStatusWaiting)
	if err != nil {
		return err
	}
	if payment != nil {
		// user already send payment
		// wait admin to accept or reject to send another payment for this course
		return ctx.Forbidden()
	}

	// calculate price
	originalPrice := course.Price
	if course.Options.Discount {
		originalPrice = course.DiscountedPrice
	}
	// TODO: calculate code

	// auto enroll if course free
	if originalPrice == 0.0 {
		enroll = &store.Enroll{
			UserID:   ctx.CurrentUserID,
			CourseID: ctx.CourseID,
		}
		err = c.db.EnrollSave(enroll)
		if err != nil {
			return err
		}
		return ctx.NoContent()
	}

	// create payment
	payment = &store.Payment{
		CourseID:      ctx.CourseID,
		UserID:        ctx.CurrentUserID,
		OriginalPrice: originalPrice,
		Price:         originalPrice,
		Code:          ctx.Payload.Code,
		URL:           ctx.Payload.URL,
	}

	err = c.db.PaymentSave(payment)
	if err != nil {
		return err
	}

	return ctx.NoContent()
}
