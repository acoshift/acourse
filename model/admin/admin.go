package admin

import "time"

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
