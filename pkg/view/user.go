package view

import (
	"github.com/acoshift/acourse/pkg/model"
)

// User type
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// UserTiny type
type UserTiny struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// UserMe type
type UserMe struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
	Role     *Role  `json:"role"`
}

// UserCollection type
type UserCollection []*User

// UserTinyCollection type
type UserTinyCollection []*UserTiny

// ToUser builds an User view from an User model
func ToUser(m *model.User) *User {
	return &User{
		ID:       m.ID,
		Username: m.Username,
		Name:     m.Name,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
	}
}

// ToUserCollection builds an UserCollection view from User models
func ToUserCollection(ms []*model.User) UserCollection {
	rs := make(UserCollection, len(ms))
	for i := range ms {
		rs[i] = ToUser(ms[i])
	}
	return rs
}

// ToUserMe builds an UserMe view from an User model
func ToUserMe(m *model.User, role *Role) *UserMe {
	return &UserMe{
		ID:       m.ID,
		Username: m.Username,
		Name:     m.Name,
		Photo:    m.Photo,
		AboutMe:  m.AboutMe,
		Role:     role,
	}
}

// ToUserTiny builds an UserTiny view from an User model
func ToUserTiny(m *model.User) *UserTiny {
	return &UserTiny{
		ID:       m.ID,
		Username: m.Username,
		Name:     m.Name,
		Photo:    m.Photo,
	}
}

// ToUserTinyCollection build an UserTinyCollection view from User models
func ToUserTinyCollection(ms []*model.User) UserTinyCollection {
	rs := make(UserTinyCollection, len(ms))
	for i := range ms {
		rs[i] = ToUserTiny(ms[i])
	}
	return rs
}
