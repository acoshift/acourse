package model

import (
	"github.com/acoshift/ds"
)

// User model
type User struct {
	ds.Model
	ds.StampModel
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

// Kind implements Kind interface
func (*User) Kind() string { return "User" }

// Users type
type Users []*User
