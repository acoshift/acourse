package user

import (
	"github.com/acoshift/acourse/pkg/acourse"
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

func toUser(x *user) *acourse.User {
	return &acourse.User{
		Id:       x.ID(),
		Username: x.Username,
		Name:     x.Name,
		Photo:    x.Photo,
		AboutMe:  x.AboutMe,
	}
}

func toUsers(xs []*user) []*acourse.User {
	rs := make([]*acourse.User, len(xs))
	for i, x := range xs {
		rs[i] = toUser(x)
	}
	return rs
}
