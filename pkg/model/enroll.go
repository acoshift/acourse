package model

import (
	"github.com/acoshift/ds"
)

// Enroll model
type Enroll struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

// Enrolls type
type Enrolls []*Enroll
