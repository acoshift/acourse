package model

import (
	"github.com/acoshift/ds"
)

// Attend model
type Attend struct {
	ds.Model
	ds.StampModel
	UserID   string
	CourseID string
}

// Kind implements Kind interface
func (*Attend) Kind() string { return "Attend" }

// Attends type
type Attends []*Attend
