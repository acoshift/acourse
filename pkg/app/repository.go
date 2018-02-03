package app

import (
	"context"

	"github.com/acoshift/acourse/pkg/entity"
)

// Repository is the app's repository
type Repository interface {
	// Auth
	StoreMagicLink(ctx context.Context, linkID string, userID string) error
	FindMagicLink(ctx context.Context, linkID string) (string, error)
	CanAcquireMagicLink(ctx context.Context, email string) (bool, error)

	// User
	SaveUser(ctx context.Context, x *entity.User) error
	GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, error)
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	GetUserFromUsername(ctx context.Context, username string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]*entity.User, error)
	CountUsers(ctx context.Context) (int64, error)
	IsUserExists(ctx context.Context, id string) (bool, error)
	CreateUser(ctx context.Context, x *entity.User) error

	// Course
	SaveCourse(ctx context.Context, x *entity.Course) error
	GetCourses(ctx context.Context, courseIDs []string) ([]*entity.Course, error)
	GetCourse(ctx context.Context, courseID string) (*entity.Course, error)
	GetCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error)
	GetCourseContent(ctx context.Context, courseContentID string) (*entity.CourseContent, error)
	GetCourseIDFromURL(ctx context.Context, url string) (string, error)
	ListCourses(ctx context.Context, limit, offset int64) ([]*entity.Course, error)
	ListPublicCourses(ctx context.Context) ([]*entity.Course, error)
	ListOwnCourses(ctx context.Context, userID string) ([]*entity.Course, error)
	ListEnrolledCourses(ctx context.Context, userID string) ([]*entity.Course, error)
	CountCourses(ctx context.Context) (int64, error)

	// Enroll
	Enroll(ctx context.Context, userID string, courseID string) error
	IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error)

	// Payment
	CreatePayment(ctx context.Context, x *entity.Payment) error
	AcceptPayment(ctx context.Context, x *entity.Payment) error
	RejectPayment(ctx context.Context, x *entity.Payment) error
	GetPayments(ctx context.Context, paymentIDs []string) ([]*entity.Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error)
	HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error)
	ListHistoryPayments(ctx context.Context, limit, offset int64) ([]*entity.Payment, error)
	ListPendingPayments(ctx context.Context, limit, offset int64) ([]*entity.Payment, error)
	CountHistoryPayments(ctx context.Context) (int64, error)
	CountPendingPayments(ctx context.Context) (int64, error)

	// Assignment
	GetAssignments(ctx context.Context, courseID string) ([]*entity.Assignment, error)
}
