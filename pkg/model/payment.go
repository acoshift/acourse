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

// PaymentStatus type
type PaymentStatus string

// Payment status
const (
	PaymentStatusWaiting  PaymentStatus = "waiting"
	PaymentStatusApproved PaymentStatus = "approved"
	PaymentStatusRejected PaymentStatus = "rejected"
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
