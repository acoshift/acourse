package app

import (
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
)

// UsersReply type
type UsersReply struct {
	Users view.Users `json:"users"`
}

// UserReply type
type UserReply struct {
	User *view.User `json:"user"`
	Role *view.Role `json:"role,omitempty"`
}

// PaymentsReply type
type PaymentsReply struct {
	Payments model.Payments
	Users    model.Users
	Courses  model.Courses
}

// Expose exposes reply
func (reply *PaymentsReply) Expose() interface{} {
	return map[string]interface{}{
		"payments": reply.Payments.Expose(),
		"users":    reply.Users.Expose(),
		"courses":  reply.Courses.Expose(),
	}
}

// CoursesReply type
type CoursesReply struct {
	Courses     model.Courses
	Users       model.Users
	EnrollCount map[string]int
}

// Expose exposes reply
func (reply *CoursesReply) Expose() interface{} {
	r := map[string]interface{}{"courses": reply.Courses.Expose()}
	if reply.Users != nil {
		r["users"] = reply.Users.ExposeMap()
	}
	if reply.EnrollCount != nil {
		r["enrollCount"] = reply.EnrollCount
	}
	return r
}
