package model

import (
	"time"
)

// Attend model
type Attend struct {
	Base
	Stampable
	UserID   string
	CourseID string
	At       time.Time `datastore:",noindex"`
}
