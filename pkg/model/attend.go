package model

import (
	"github.com/acoshift/ds"
)

// Attend model
type Attend struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

// Attends type
type Attends []*Attend
