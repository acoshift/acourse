package model

// Assignment model
type Assignment struct {
	Base
	Stampable
	UserID             string
	CourseID           string
	CourseAssignmentID string `datastore:",noindex"`
	URL                string `datastore:",noindex"`
}
