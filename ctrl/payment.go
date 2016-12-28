package ctrl

import (
	"acourse/app"
	"acourse/model"
	"acourse/store"
	"acourse/view"
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
	role, err := c.db.RoleFindByUserID(ctx.CurrentUserID)
	if err != nil {
		return err
	}

	// only admin can access
	if !role.Admin {
		return ctx.Forbidden()
	}

	xs, err := c.db.PaymentList(model.PaymentStatusWaiting)
	if err != nil {
		return err
	}

	res := make(view.PaymentCollection, len(xs))
	for i, x := range xs {
		user, err := c.db.UserMustGet(x.UserID)
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

// Approve runs approve action
func (c *PaymentController) Approve(ctx *app.PaymentApproveContext) error {
	role, err := c.db.RoleFindByUserID(ctx.CurrentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return ctx.Forbidden()
	}

	payment, err := c.db.PaymentGet(ctx.PaymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return ctx.NotFound()
	}
	payment.Status = model.PaymentStatusApproved

	// Add user to enroll
	enroll := &model.Enroll{
		UserID:   payment.UserID,
		CourseID: payment.CourseID,
	}
	err = c.db.EnrollSave(enroll)
	if err != nil {
		return err
	}

	err = c.db.PaymentSave(payment)
	if err != nil {
		return err
	}

	return ctx.OK()
}

// Reject runs reject action
func (c *PaymentController) Reject(ctx *app.PaymentRejectContext) error {
	role, err := c.db.RoleFindByUserID(ctx.CurrentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return ctx.Forbidden()
	}

	payment, err := c.db.PaymentGet(ctx.PaymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return ctx.NotFound()
	}
	payment.Status = model.PaymentStatusRejected

	err = c.db.PaymentSave(payment)
	if err != nil {
		return err
	}

	return ctx.OK()
}
