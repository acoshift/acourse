package acourse

import (
	"github.com/acoshift/acourse/pkg/model"
)

// ToUser builds an User message from an User model
func ToUser(x *model.User) *User {
	return &User{
		Id:       x.ID,
		Username: x.Username,
		Name:     x.Name,
		Photo:    x.Photo,
		AboutMe:  x.AboutMe,
	}
}

// ToUsers builds repeated User message from User models
func ToUsers(xs model.Users) []*User {
	rs := make([]*User, len(xs))
	for i, x := range xs {
		rs[i] = ToUser(x)
	}
	return rs
}

// ToRole builds a Role message from Role model
func ToRole(x *model.Role) *Role {
	return &Role{
		Admin:      x.Admin,
		Instructor: x.Instructor,
	}
}
