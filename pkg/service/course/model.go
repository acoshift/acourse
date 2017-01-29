package course

import (
	"github.com/acoshift/ds"
)

type attend struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

const (
	kindAttend = "Attend"
)
