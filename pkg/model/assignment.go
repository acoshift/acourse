package model

import (
	"github.com/acoshift/ds"
)

// Assignment model
type Assignment struct {
	ds.Model
	ds.StampModel
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

// Kind implements Kind interface
func (*Assignment) Kind() string { return "Assignment" }

// Assignments type
type Assignments []*Assignment

// UserAssignment model
type UserAssignment struct {
	ds.Model
	ds.StampModel
	AssignmentID string
	UserID       string
	URL          string `datastore:",noindex"`
}

// Kind implements Kind interface
func (*UserAssignment) Kind() string { return "UserAssignment" }

// UserAssignments type
type UserAssignments []*UserAssignment
