package model

import (
	"github.com/acoshift/ds"
)

// Enroll model
type Enroll struct {
	ds.Model
	ds.StampModel
	UserID   string
	CourseID string
}

// Kind implements Kind interface
func (*Enroll) Kind() string { return "Enroll" }

// Enrolls type
type Enrolls []*Enroll
