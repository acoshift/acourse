package app

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

// Repository is the app's repository
type Repository interface {
	// Auth
	StoreMagicLink(ctx context.Context, linkID string, userID string) error

	// User
	SaveUser(ctx context.Context, x *User) error
	GetUsers(ctx context.Context, userIDs []string) ([]*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUserFromUsername(ctx context.Context, username string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]*User, error)
	CountUsers(ctx context.Context) (int64, error)
	IsUserExists(ctx context.Context, id string) (bool, error)
	CreateUser(ctx context.Context, x *User) error

	// Course
	SaveCourse(ctx context.Context, x *Course) error
	GetCourses(ctx context.Context, courseIDs []string) ([]*Course, error)
	GetCourse(ctx context.Context, courseID string) (*Course, error)
	GetCourseContents(ctx context.Context, courseID string) ([]*CourseContent, error)
	GetCourseContent(ctx context.Context, courseContentID string) (*CourseContent, error)
	GetCourseIDFromURL(ctx context.Context, url string) (string, error)
	ListCourses(ctx context.Context, limit, offset int64) ([]*Course, error)
	ListPublicCourses(ctx context.Context, cachePool *redis.Pool, cachePrefix string) ([]*Course, error)
	ListOwnCourses(ctx context.Context, userID string) ([]*Course, error)
	ListEnrolledCourses(ctx context.Context, userID string) ([]*Course, error)
	CountCourses(ctx context.Context) (int64, error)

	// Enroll
	Enroll(ctx context.Context, userID string, courseID string) error
	IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error)

	// Payment
	CreatePayment(ctx context.Context, x *Payment) error
	AcceptPayment(ctx context.Context, x *Payment) error
	RejectPayment(ctx context.Context, x *Payment) error
	GetPayments(ctx context.Context, paymentIDs []string) ([]*Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error)
	ListHistoryPayments(ctx context.Context, limit, offset int64) ([]*Payment, error)
	ListPendingPayments(ctx context.Context, limit, offset int64) ([]*Payment, error)
	CountHistoryPayments(ctx context.Context) (int64, error)
	CountPendingPayments(ctx context.Context) (int64, error)

	// Assignment
	GetAssignments(ctx context.Context, courseID string) ([]*Assignment, error)
}
