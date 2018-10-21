package payment

// Accept command
type Accept struct {
	ID string
}

// Reject command
type Reject struct {
	ID      string
	Message string
}

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
