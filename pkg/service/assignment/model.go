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

func (x *assignment) NewKey() {
	x.NewIncomplateKey(kindAssignment, nil)
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

type userAssignment struct {
	ds.StringIDModel
	ds.StampModel
	AssignmentID string
	UserID       string
	URL          string `datastore:",noindex"`
}

func (x *userAssignment) NewKey() {
	x.NewIncomplateKey(kindUserAssignment, nil)
}

func toUserAssignment(x *userAssignment) *acourse.UserAssignment {
	return &acourse.UserAssignment{
		Id:           x.ID(),
		CreatedAt:    x.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    x.UpdatedAt.Format(time.RFC3339),
		AssignmentId: x.AssignmentID,
		UserId:       x.UserID,
		Url:          x.URL,
	}
}

func toUserAssignments(xs []*userAssignment) []*acourse.UserAssignment {
	rs := make([]*acourse.UserAssignment, len(xs))
	for i, x := range xs {
		rs[i] = toUserAssignment(x)
	}
	return rs
}
