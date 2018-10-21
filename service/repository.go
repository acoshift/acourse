package service

import (
	"context"

	"github.com/acoshift/acourse/model/course"
)

// Repository is the service storage
type Repository interface {
	RegisterCourse(ctx context.Context, x *RegisterCourse) (courseID string, err error)
	GetCourse(ctx context.Context, courseID string) (*course.Course, error)
	UpdateCourse(ctx context.Context, x *UpdateCourseModel) error

	RegisterPayment(ctx context.Context, x *RegisterPayment) error
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)

	RegisterEnroll(ctx context.Context, userID string, courseID string) error
}
