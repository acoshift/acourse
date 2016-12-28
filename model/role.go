package model

// Role model store user's role
type Role struct {
	Base
	Stampable

	// roles
	Admin      bool
	Instructor bool
}
