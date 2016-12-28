package model

import (
	"time"
)

// Stampable is the database timestamp model
type Stampable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Stamp stamps current time to model
func (m *Stampable) Stamp() {
	m.UpdatedAt = time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = m.UpdatedAt
	}
}
