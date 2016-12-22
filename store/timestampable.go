package store

import "time"

// Timestampable is the database timestamp model
type Timestampable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Stamp current time to model
func (m *Timestampable) Stamp() {
	m.UpdatedAt = time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = m.UpdatedAt
	}
}
