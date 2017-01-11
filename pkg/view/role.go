package view

import (
	"github.com/acoshift/acourse/pkg/model"
)

// Role type
type Role struct {
	Admin      bool `json:"admin"`
	Instructor bool `json:"instructor"`
}

// ToRole builds a Role view fromn a Role model
func ToRole(m *model.Role) *Role {
	return &Role{
		Admin:      m.Admin,
		Instructor: m.Instructor,
	}
}
