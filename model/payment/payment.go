package payment

// SetStatus sets payment status
type SetStatus struct {
	ID     string
	Status int
}

// HasPending checks is has course pending payment
type HasPending struct {
	UserID   string
	CourseID string

	Result bool
}
