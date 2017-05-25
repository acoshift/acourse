package ds

import (
	"time"
)

// StampModel is the stampable model
type StampModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Stamp stamps current time to model
func (x *StampModel) Stamp() {
	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
	}
}

// Stampable interface
type Stampable interface {
	Stamp()
}
