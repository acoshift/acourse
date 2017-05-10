package view

import (
	"github.com/acoshift/acourse/pkg/model"
)

// Page type provides layout data like title, description, and og
type Page struct {
	Title string
	Desc  string
	Image string
	URL   string
}

// IndexData type
type IndexData struct {
	*Page
	Courses []*model.Course
}

// AuthData type
type AuthData struct {
	*Page
}
