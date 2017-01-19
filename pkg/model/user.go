package model

import (
	"github.com/acoshift/ds"
)

// User model
type User struct {
	ds.StringIDModel
	ds.StampModel
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

// Users type
type Users []*User
