package user

import (
	"github.com/acoshift/ds"
)

type user struct {
	ds.StringIDModel
	ds.StampModel
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

type role struct {
	ds.StringIDModel
	ds.StampModel

	// roles
	Admin      bool
	Instructor bool
}

const (
	kindUser = "User"
	kindRole = "Role"
)
