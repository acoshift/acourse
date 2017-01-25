package model

import (
	"github.com/acoshift/ds"
)

// UserAssignment model
type UserAssignment struct {
	ds.StringIDModel
	ds.StampModel
	AssignmentID string
	UserID       string
	URL          string `datastore:",noindex"`
}

// UserAssignments type
type UserAssignments []*UserAssignment
