package model

import (
	"github.com/acoshift/ds"
)

// Favorite model
type Favorite struct {
	ds.Model
	ds.StampModel
	UserID   string
	CourseID string
}

// Kind implements Kind interface
func (*Favorite) Kind() string { return "Favorite" }
