package app

import (
	"context"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/payload"
)

// HealthService interface
type HealthService interface {
	Check(context.Context) error
}

// UserService interface
type UserService interface {
	GetUsers(context.Context, *IDsRequest) (*UsersReply, error)
	GetMe(context.Context) (*UserReply, error)
	UpdateMe(context.Context, *UserRequest) error
}

// PaymentService interface
type PaymentService interface {
	ListPayments(context.Context, *PaymentListRequest) (*PaymentsReply, error)
	ApprovePayments(context.Context, *IDsRequest) error
	RejectPayments(context.Context, *IDsRequest) error
}

// CourseService interface
type CourseService interface {
	ListCourses(context.Context, *CourseListRequest) (*CoursesReply, error)
}

// EmailService interface
type EmailService interface {
	SendEmail(context.Context, *EmailRequest) error
}

// IDsRequest type
type IDsRequest struct {
	IDs []string `json:"ids"`
}

// UniqueIDs filters out duplicated ID from IDs
func (req *IDsRequest) UniqueIDs() []string {
	idMap := map[string]bool{}
	for _, id := range req.IDs {
		idMap[id] = true
	}
	res := make([]string, len(idMap))
	for id := range idMap {
		res = append(res, id)
	}
	return res
}

// UserRequest type
type UserRequest struct {
	User *payload.RawUser `json:"user"`
}

// PaymentListRequest type
type PaymentListRequest struct {
	Offset  *int  `json:"offset"`
	Limit   *int  `json:"limit"`
	History *bool `json:"history"`
}

// EmailRequest type
type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// CourseListRequest type
type CourseListRequest struct {
	Public  *bool `json:"public"`
	Student bool  `json:"student"`
}

// UsersReply type
type UsersReply struct {
	Users model.Users
}

// Expose exposes reply
func (reply *UsersReply) Expose() interface{} {
	return map[string]interface{}{
		"users": reply.Users.Expose(),
	}
}

// UserReply type
type UserReply struct {
	User *model.User
	Role *model.Role
}

// Expose exposes reply
func (reply *UserReply) Expose() interface{} {
	r := map[string]interface{}{
		"user": reply.User.Expose(),
	}
	if reply.Role != nil {
		r["role"] = reply.Role.Expose()
	}
	return r
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
	Courses  model.Courses
	Users    model.Users
	Students map[string]int
}

// Expose exposes reply
func (reply *CoursesReply) Expose() interface{} {
	r := map[string]interface{}{
		"courses": reply.Courses.Expose(),
		"users":   reply.Users.ExposeMap(),
	}
	if reply.Students != nil {
		r["students"] = reply.Students
	}
	return r
}
