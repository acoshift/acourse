package model

// Assignment model
type Assignment struct {
	Base
	Stampable
	CourseID    string
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

// Assignments type
type Assignments []*Assignment

// UserAssignment model
type UserAssignment struct {
	Base
	Stampable
	AssignmentID string
	URL          string `datastore:",noindex"`
}

// UserAssignments type
type UserAssignments []*UserAssignment
