package model

import (
	"github.com/acoshift/ds"
)

// Assignment model
type Assignment struct {
	ds.StringIDModel
	ds.StampModel
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

// Assignments type
type Assignments []*Assignment

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
