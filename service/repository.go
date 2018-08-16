package service

import (
	"context"

	"github.com/acoshift/acourse/entity"
)

// Repository is the service storage
type Repository interface {
	StoreMagicLink(ctx context.Context, linkID string, userID string) error
	FindMagicLink(ctx context.Context, linkID string) (string, error)
	CanAcquireMagicLink(ctx context.Context, email string) (bool, error)

	GetEmailSignInUserByEmail(ctx context.Context, email string) (*entity.EmailSignInUser, error)

	RegisterUser(ctx context.Context, x *RegisterUser) error
	UpdateUser(ctx context.Context, x *UpdateUser) error
	SetUserImage(ctx context.Context, userID string, image string) error
	IsUserExists(ctx context.Context, userID string) (exists bool, err error)

	RegisterCourse(ctx context.Context, x *RegisterCourse) (courseID string, err error)
	GetCourse(ctx context.Context, courseID string) (*entity.Course, error)
	UpdateCourse(ctx context.Context, x *UpdateCourseModel) error
	SetCourseImage(ctx context.Context, courseID string, image string) error
	SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error

	RegisterPayment(ctx context.Context, x *RegisterPayment) error
	GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error)
	SetPaymentStatus(ctx context.Context, paymentID string, status int) error
	HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error)

	RegisterEnroll(ctx context.Context, userID string, courseID string) error
	IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error)
}
