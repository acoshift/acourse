package view

import (
	"github.com/acoshift/acourse/pkg/model"
)

// User view
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// ToUser builds User view from User model
func ToUser(x *model.User) *User {
	return &User{x.ID, x.Username, x.Name, x.Photo, x.AboutMe}
}

// Users view
type Users []*User

// ToUsers builds Users view from Users model
func ToUsers(xs model.Users) Users {
	rs := make(Users, len(xs))
	for i, x := range xs {
		rs[i] = ToUser(x)
	}
	return rs
}

// UserMap view
type UserMap map[string]*User

// ToUserMap builds User map from Users model
func ToUserMap(xs model.Users) UserMap {
	m := make(UserMap, len(xs))
	for _, x := range xs {
		m[x.ID] = ToUser(x)
	}
	return m
}

// UserTiny view
type UserTiny struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// ToUserTiny builds User tiny view from User model
func ToUserTiny(x *model.User) *UserTiny {
	return &UserTiny{x.ID, x.Username, x.Name, x.Photo}
}

// UsersTiny view
type UsersTiny []*UserTiny

// ToUsersTiny builds Users tiny view from Users model
func ToUsersTiny(xs model.Users) UsersTiny {
	rs := make(UsersTiny, len(xs))
	for i, x := range xs {
		rs[i] = ToUserTiny(x)
	}
	return rs
}

// UserTinyMap view
type UserTinyMap map[string]*UserTiny

// ToUserTinyMap builds User tiny map from Users model
func ToUserTinyMap(xs model.Users) UserTinyMap {
	m := make(UserTinyMap, len(xs))
	for _, x := range xs {
		m[x.ID] = ToUserTiny(x)
	}
	return m
}
