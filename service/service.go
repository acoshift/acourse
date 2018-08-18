package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/acoshift/go-firebase-admin"

	"github.com/acoshift/acourse/email"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/notify"
)

// Config is service config
type Config struct {
	Repository         Repository
	Auth               *firebase.Auth
	EmailSender        email.Sender
	BaseURL            string
	FileStorage        file.Storage
	ImageResizeEncoder image.JPEGResizeEncoder
	AdminNotifier      notify.AdminNotifier
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
	EnrollCourse(ctx context.Context, courseID string, price float64, paymentImage *multipart.FileHeader) error

	UpdateProfile(ctx context.Context, x *Profile) error

	CreateCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error)
	GetCourseContent(ctx context.Context, contentID string) (*entity.CourseContent, error)
	UpdateCourseContent(ctx context.Context, contentID string, title string, desc string, videoID string) error
	DeleteCourseContent(ctx context.Context, contentID string) error
}

// New creates new service
func New(cfg Config) Service {
	return &svc{cfg}
}

type svc struct {
	Config
}
