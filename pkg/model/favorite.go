package model

import (
	"github.com/acoshift/ds"
)

// Favorite model
type Favorite struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}
