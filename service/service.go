package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/entity"
)

// Config is service config
type Config struct {
	Repository     Repository
	BaseURL        string
	Location       *time.Location
	OpenIDCallback string
}

// Service type
type Service interface {
	GenerateOpenIDURI(ctx context.Context, provider string) (redirectURI string, state string, err error)
	SignInOpenIDCallback(ctx context.Context, uri string, state string) (userID string, err error)

	CreateCourse(ctx context.Context, x *CreateCourse) (courseID string, err error)
	UpdateCourse(ctx context.Context, x *UpdateCourse) error
	EnrollCourse(ctx context.Context, courseID string, price float64, paymentImage *multipart.FileHeader) error

	CreateCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error)
	GetCourseContent(ctx context.Context, contentID string) (*entity.CourseContent, error)
	ListCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error)
	UpdateCourseContent(ctx context.Context, contentID string, title string, desc string, videoID string) error
	DeleteCourseContent(ctx context.Context, contentID string) error
}

// New creates new service
func New(cfg Config) Service {
	s := &svc{cfg}

	dispatcher.Register(s.signUp)
	dispatcher.Register(s.sendPasswordResetEmail)
	dispatcher.Register(s.signInPassword)
	dispatcher.Register(s.updateProfile)
	dispatcher.Register(s.acceptPayment)
	dispatcher.Register(s.rejectPayment)
	return s
}

type svc struct {
	Config
}
