package store

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/ds"
)

// PaymentListOptions type
type PaymentListOptions struct {
	Status *model.PaymentStatus
}

// PaymentListOption type
type PaymentListOption func(*PaymentListOptions)

// PaymentListOptionStatus sets status to options
func PaymentListOptionStatus(status model.PaymentStatus) PaymentListOption {
	return func(args *PaymentListOptions) {
		args.Status = &status
	}
}

// PaymentList list all payments
func (c *DB) PaymentList(ctx context.Context, opts ...PaymentListOption) (model.Payments, error) {
	var xs []*model.Payment

	opt := &PaymentListOptions{}
	for _, setter := range opts {
		setter(opt)
	}

	qs := []ds.Query{}

	if opt.Status != nil {
		qs = append(qs, ds.Filter("Status =", string(*opt.Status)))
	}

	qs = append(qs, ds.Order("CreatedAt"))

	err := c.client.Query(ctx, &model.Payment{}, &xs, qs...)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return xs, nil
}

// PaymentSave saves a payment to database
func (c *DB) PaymentSave(ctx context.Context, x *model.Payment) error {
	return c.client.Save(ctx, x)
}

// PaymentSaveMulti saves multiple payments to database
func (c *DB) PaymentSaveMulti(ctx context.Context, payments model.Payments) error {
	err := c.client.SaveMulti(ctx, payments)
	if err != nil {
		return err
	}
	return nil
}

// PaymentGet retrieves a payment from database
func (c *DB) PaymentGet(ctx context.Context, paymentID string) (*model.Payment, error) {
	var x model.Payment
	err := c.client.GetByID(ctx, paymentID, &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// PaymentFind finds a payment with user id and course id
func (c *DB) PaymentFind(ctx context.Context, userID, courseID string, status model.PaymentStatus) (*model.Payment, error) {
	var x model.Payment
	err := c.client.QueryFirst(ctx, &x,
		ds.Filter("UserID =", userID),
		ds.Filter("CourseID =", courseID),
		ds.Filter("Status =", string(status)),
	)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// PaymentGetMulti retrieves multiple payments from database
func (c *DB) PaymentGetMulti(ctx context.Context, paymentIDs []string) (model.Payments, error) {
	if len(paymentIDs) == 0 {
		return []*model.Payment{}, nil
	}

	payments := make([]*model.Payment, len(paymentIDs))
	err := c.client.GetByIDs(ctx, paymentIDs, &model.Payment{}, &payments)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
