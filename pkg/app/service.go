package app

import (
	"context"

	"github.com/acoshift/acourse/pkg/payload"
	"github.com/acoshift/acourse/pkg/view"
)

// UserService interface
type UserService interface {
	GetUsers(context.Context, *IDsRequest) (*UsersReply, error)
	GetMe(context.Context) (*UserMeReply, error)
	UpdateMe(context.Context, *UserRequest) error
}

// IDsRequest type
type IDsRequest struct {
	IDs []string `json:"ids"`
}

// UserRequest type
type UserRequest struct {
	User *payload.RawUser `json:"user"`
}

// UsersReply type
type UsersReply struct {
	Users view.UserCollection `json:"users"`
}

// UserMeReply type
type UserMeReply struct {
	User *view.UserMe `json:"user"`
}
