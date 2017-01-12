package model

// Role model store user's role
type Role struct {
	Base
	Stampable

	// roles
	Admin      bool
	Instructor bool
}

// Expose exposes model
func (x *Role) Expose() interface{} {
	return map[string]interface{}{
		"admin":      x.Admin,
		"instructor": x.Instructor,
	}
}
