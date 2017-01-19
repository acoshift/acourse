package model

import (
	"github.com/acoshift/ds"
)

// Role model store user's role
type Role struct {
	ds.StringIDModel
	ds.StampModel

	// roles
	Admin      bool
	Instructor bool
}
