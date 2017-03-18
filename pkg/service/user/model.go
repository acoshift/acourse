package user

import (
	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
)

type userModel struct {
	ds.StringIDModel
	ds.StampModel
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

func (x *userModel) NewKey() {
	x.NewIncomplateKey(kindUser, nil)
}

type roleModel struct {
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

func toUser(x *userModel) *acourse.User {
	return &acourse.User{
		Id:       x.ID(),
		Username: x.Username,
		Name:     x.Name,
		Photo:    x.Photo,
		AboutMe:  x.AboutMe,
	}
}

func toUsers(xs []*userModel) []*acourse.User {
	rs := make([]*acourse.User, len(xs))
	for i, x := range xs {
		rs[i] = toUser(x)
	}
	return rs
}
