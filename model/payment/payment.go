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
