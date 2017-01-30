package course

import (
	"github.com/acoshift/ds"
)

type attendModel struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

type enrollModel struct {
	ds.StringIDModel
	ds.StampModel
	UserID   string
	CourseID string
}

const (
	kindAttend = "Attend"
	kindEnroll = "Enroll"
)
