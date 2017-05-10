package model

import "time"

// Payment model
type Payment struct {
	id        string
	UserID    string
	CourseID  string
	Image     string
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentStatus type
type PaymentStatus int

// PaymentStatus values
const (
	Pending PaymentStatus = iota
	Accepted
	Rejected
)

// ID returns payment id
func (x *Payment) ID() string {
	return x.id
}
