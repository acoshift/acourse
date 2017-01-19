package model

import (
	"github.com/acoshift/ds"
)

// Role model store user's role
type Role struct {
	ds.Model
	ds.StampModel

	// roles
	Admin      bool
	Instructor bool
}

// Kind implements Kind interface
func (*Role) Kind() string { return "Role" }
