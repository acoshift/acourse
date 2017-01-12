package model

import (
	"time"
)

// Payment model
type Payment struct {
	Base
	Stampable
	UserID        string
	CourseID      string
	OriginalPrice float64 `datastore:",noindex"`
	Price         float64 `datastore:",noindex"`
	Code          string
	URL           string `datastore:",noindex"`
	Status        PaymentStatus
	At            time.Time
}

// Payments type
type Payments []*Payment

// PaymentStatus type
type PaymentStatus string

// Payment status
const (
	PaymentStatusWaiting  PaymentStatus = "waiting"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusRejected PaymentStatus = "rejected"
)

// PaymentView type
type PaymentView int

// PaymentView
const (
	PaymentViewDefault PaymentView = iota
)

// Approve approves a payment
func (x *Payment) Approve() {
	x.Status = PaymentStatusApproved
	x.At = time.Now()
}

// Reject rejects a payment
func (x *Payment) Reject() {
	x.Status = PaymentStatusRejected
	x.At = time.Now()
}

// SetView sets view to model
func (x *Payment) SetView(v PaymentView) {
	x.view = v
}

// SetView sets view to model
func (xs Payments) SetView(v PaymentView) {
	for _, x := range xs {
		x.SetView(v)
	}
}

// Expose exposes model
func (x *Payment) Expose() interface{} {
	if x.view == nil {
		return nil
	}
	switch x.view.(PaymentView) {
	case PaymentViewDefault:
		return map[string]interface{}{
			"id":            x.ID,
			"userId":        x.UserID,
			"courseId":      x.CourseID,
			"originalPrice": x.OriginalPrice,
			"price":         x.Price,
			"code":          x.Code,
			"url":           x.URL,
			"status":        string(x.Status),
			"createdAt":     x.CreatedAt,
			"updatedAt":     x.UpdatedAt,
			"at":            x.At,
		}
	default:
		return nil
	}
}

// Expose exposes model
func (xs Payments) Expose() interface{} {
	rs := make([]interface{}, len(xs))
	for i, x := range xs {
		rs[i] = x.Expose()
	}
	return rs
}
