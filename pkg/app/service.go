package app

import (
	"context"

	"github.com/acoshift/acourse/pkg/payload"
	"github.com/acoshift/acourse/pkg/view"
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

// UsersReply type
type UsersReply struct {
	Users view.UserCollection `json:"users"`
}

// UserReply type
type UserReply struct {
	User *view.User `json:"user"`
	Role *view.Role `json:"role,omitempty"`
}

// PaymentsReply type
type PaymentsReply struct {
	Payments view.PaymentCollection    `json:"payments"`
	Users    view.UserTinyCollection   `json:"users"`
	Courses  view.CourseMiniCollection `json:"courses"`
}

// EmailRequest type
type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
