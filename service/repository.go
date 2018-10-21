package service

import (
	"context"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
)

// Repository is the service storage
type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	RegisterCourse(ctx context.Context, x *RegisterCourse) (courseID string, err error)
	GetCourse(ctx context.Context, courseID string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, x *UpdateCourseModel) error
	SetCourseImage(ctx context.Context, courseID string, image string) error
	SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error

	RegisterCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error)
	GetCourseContent(ctx context.Context, contentID string) (*course.Content, error)
	ListCourseContents(ctx context.Context, courseID string) ([]*course.Content, error)
	UpdateCourseContent(ctx context.Context, contentID, title, desc, videoID string) error
	DeleteCourseContent(ctx context.Context, contentID string) error

	RegisterPayment(ctx context.Context, x *RegisterPayment) error
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	SetPaymentStatus(ctx context.Context, paymentID string, status int) error
	HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error)

	RegisterEnroll(ctx context.Context, userID string, courseID string) error
	IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error)
}
