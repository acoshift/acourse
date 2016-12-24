package ctrl

import (
	"acourse/app"
	"acourse/store"
)

// PaymentController implements PaymentController interface
type PaymentController struct {
	db *store.DB
}

// NewPaymentController creates new controller
func NewPaymentController(db *store.DB) *PaymentController {
	return &PaymentController{db}
}

// List runs list action
func (c *PaymentController) List(ctx *app.PaymentListContext) error {
	role, err := c.db.RoleGet(ctx.CurrentUserID)
	if err != nil {
		return err
	}

	// only admin can access
	if !role.Admin {
		return ctx.Forbidded()
	}

	xs, err := c.db.PaymentList(store.PaymentStatusWaiting)
	if err != nil {
		return err
	}

	res := make(app.PaymentCollectionView, len(xs))
	for i, x := range xs {
		user, err := c.db.UserGet(x.UserID)
		if err != nil {
			return err
		}
		course, err := c.db.CourseGet(x.CourseID)
		if err != nil {
			return err
		}
		res[i] = ToPaymentView(x, ToUserTinyView(user), ToCourseMiniView(course))
	}

	return ctx.OK(res)
}
