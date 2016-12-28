package model

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
}

// PaymentStatus type
type PaymentStatus string

// Payment status
const (
	PaymentStatusWaiting  = "waiting"
	PaymentStatusApproved = "approved"
	PaymentStatusRejected = "rejected"
)
