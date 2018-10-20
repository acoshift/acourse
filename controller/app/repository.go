package app

import (
	"context"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
)

// Repository is app storage
type Repository interface {
	GetCourse(ctx context.Context, courseID string) (*Course, error)
	GetCourseIDByURL(ctx context.Context, url string) (string, error)
	IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error)
	HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error)
	GetCourseContents(ctx context.Context, courseID string) ([]*course.Content, error)
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	FindAssignmentsByCourseID(ctx context.Context, courseID string) ([]*entity.Assignment, error)
	ListPublicCourses(ctx context.Context) ([]*PublicCourse, error)
	ListOwnCourses(ctx context.Context, userID string) ([]*OwnCourse, error)
	ListEnrolledCourses(ctx context.Context, userID string) ([]*EnrolledCourse, error)
}
