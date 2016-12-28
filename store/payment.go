package store

import (
	"acourse/model"

	"cloud.google.com/go/datastore"
)

const kindPayment = "Payment"

// PaymentList list all payments
func (c *DB) PaymentList(status model.PaymentStatus) ([]*model.Payment, error) {
	ctx, cancel := getContext()
	defer cancel()

	var xs []*model.Payment
	q := datastore.
		NewQuery(kindPayment).
		Filter("Status =", string(status))

	keys, err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}
	for i, x := range xs {
		x.SetKey(keys[i])
	}
	return xs, nil
}

// PaymentSave saves a payment to database
func (c *DB) PaymentSave(x *model.Payment) error {
	ctx, cancel := getContext()
	defer cancel()

	x.Stamp()
	if x.Key() == nil {
		x.SetKey(datastore.IncompleteKey(kindPayment, nil))
	}

	key, err := c.client.Put(ctx, x.Key(), x)
	if err != nil {
		return err
	}
	x.SetKey(key)
	return nil
}

// PaymentGet retrieves a payment from database
func (c *DB) PaymentGet(paymentID string) (*model.Payment, error) {
	id := idInt(paymentID)
	if id == 0 {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	var x model.Payment
	err := c.get(ctx, datastore.IDKey(kindPayment, id, nil), &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

// PaymentFind finds a payment with user id and course id
func (c *DB) PaymentFind(userID, courseID string, status model.PaymentStatus) (*model.Payment, error) {
	ctx, cancel := getContext()
	defer cancel()

	q := datastore.
		NewQuery(kindPayment).
		Filter("UserID =", userID).
		Filter("CourseID =", courseID).
		Filter("Status =", string(status)).
		Limit(1)

	var x model.Payment
	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}
