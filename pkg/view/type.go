package view

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/flash"
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
	flash.Flash
}

// ProfileData type
type ProfileData struct {
	*Page
	flash.Flash
}
