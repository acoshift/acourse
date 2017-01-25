package assignment

import (
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
)

const (
	kindAssignment     = "Assignment"
	kindUserAssignment = "UserAssignment"
)

type assignment struct {
	ds.StringIDModel
	ds.StampModel
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

func toAssignment(x *assignment) *acourse.Assignment {
	return &acourse.Assignment{
		Id:          x.ID(),
		CreatedAt:   x.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   x.UpdatedAt.Format(time.RFC3339),
		Title:       x.Title,
		Description: x.Description,
		Open:        x.Open,
	}
}

func toAssignments(xs []*assignment) []*acourse.Assignment {
	rs := make([]*acourse.Assignment, len(xs))
	for i, x := range xs {
		rs[i] = toAssignment(x)
	}
	return rs
}
