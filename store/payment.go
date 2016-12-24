package store

import (
	"cloud.google.com/go/datastore"
)

// Payment model
type Payment struct {
	Base
	Timestampable
	UserID        string
	CourseID      string
	OriginalPrice float64 `datastore:",noindex"`
	Price         float64 `datastore:",noindex"`
	Code          string
	URL           string `datastore:",noindex"`
	Status        PaymentStatus
}

// PaymentStatus type
type PaymentStatus string

// Payment status
const (
	PaymentStatusWaiting  = "waiting"
	PaymentStatusApproved = "approved"
	PaymentStatusRejected = "rejected"
)

const kindPayment = "Payment"

// PaymentList list all payments
func (c *DB) PaymentList(status PaymentStatus) ([]*Payment, error) {
	ctx, cancel := getContext()
	defer cancel()

	var xs []*Payment
	q := datastore.
		NewQuery(kindPayment).
		Filter("Status =", status)

	keys, err := c.getAll(ctx, q, &xs)
	if err != nil {
		return nil, err
	}
	for i, x := range xs {
		x.setKey(keys[i])
	}
	return xs, nil
}

// PaymentSave saves a payment to database
func (c *DB) PaymentSave(x *Payment) error {
	ctx, cancel := getContext()
	defer cancel()

	x.Stamp()
	if x.key == nil {
		x.setKey(datastore.IncompleteKey(kindPayment, nil))
	}

	key, err := c.client.Put(ctx, x.key, x)
	if err != nil {
		return err
	}
	x.setKey(key)
	return nil
}

// PaymentGet retrieves a payment from database
func (c *DB) PaymentGet(paymentID string) (*Payment, error) {
	id := idInt(paymentID)
	if id == 0 {
		return nil, ErrInvalidID
	}

	ctx, cancel := getContext()
	defer cancel()

	var x Payment
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
func (c *DB) PaymentFind(userID, courseID string, status PaymentStatus) (*Payment, error) {
	ctx, cancel := getContext()
	defer cancel()

	q := datastore.
		NewQuery(kindPayment).
		Filter("UserID =", userID).
		Filter("CourseID =", courseID).
		Filter("Status =", string(status)).
		Limit(1)

	var x Payment
	err := c.findFirst(ctx, q, &x)
	if notFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}
