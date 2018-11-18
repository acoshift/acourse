package admin

import "time"

// ListUsers command
type ListUsers struct {
	Limit  int64
	Offset int64

	Result []*UserItem
}

// CountUsers command
type CountUsers struct {
	Result int64
}

// ListCourses command
type ListCourses struct {
	Limit  int64
	Offset int64

	Result []*CourseItem
}

// CountCourses command
type CountCourses struct {
	Result int64
}

// GetPayment command
type GetPayment struct {
	PaymentID string

	Result Payment
}

// ListPayments command
type ListPayments struct {
	Status []int
	Limit  int64
	Offset int64

	Result []*Payment
}

// CountPayments command
type CountPayments struct {
	Status []int

	Result int64
}

// AcceptPayment command
type AcceptPayment struct {
	ID       string
	Location *time.Location
}

// RejectPayment command
type RejectPayment struct {
	ID      string
	Message string
}
