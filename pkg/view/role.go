package view

import (
	"github.com/acoshift/acourse/pkg/model"
)

// Role view
type Role struct {
	Admin      bool `json:"admin"`
	Instructor bool `json:"instructor"`
}

// ToRole builds Role view from Role model
func ToRole(x *model.Role) *Role {
	return &Role{x.Admin, x.Instructor}
}
