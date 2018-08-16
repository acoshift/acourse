package service

import (
	"context"
	"time"

	"github.com/acoshift/go-firebase-admin"

	"github.com/acoshift/acourse/email"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
)

// Config is service config
type Config struct {
	Auth               *firebase.Auth
	EmailSender        email.Sender
	BaseURL            string
	FileStorage        file.Storage
	ImageResizeEncoder image.JPEGResizeEncoder
	Location           *time.Location
	MagicLinkCallback  string
	OpenIDCallback     string
}

// Service type
type Service interface {
	SignUp(ctx context.Context, email, password string) (userID string, err error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	SendSignInMagicLinkEmail(ctx context.Context, email string) error
	GenerateOpenIDURI(ctx context.Context, provider string) (redirectURI string, state string, err error)
	SignInPassword(ctx context.Context, email, password string) (userID string, err error)
	SignInMagicLink(ctx context.Context, link string) (userID string, err error)
	SignInOpenIDCallback(ctx context.Context, uri string, state string) (userID string, err error)

	AcceptPayment(ctx context.Context, paymentID string) error
	RejectPayment(ctx context.Context, paymentID string, msg string) error

	CreateCourse(ctx context.Context, x *CreateCourse) (courseID string, err error)
	UpdateCourse(ctx context.Context, x *UpdateCourse) error
}

// New creates new service
func New(cfg Config) Service {
	return &svc{cfg}
}

type svc struct {
	Config
}
