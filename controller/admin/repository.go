package admin

import (
	"context"
)

// Repository is the admin storage
type Repository interface {
	ListUsers(ctx context.Context, limit, offset int64) ([]*UserItem, error)
	CountUsers(ctx context.Context) (int64, error)
	ListCourses(ctx context.Context, limit, offset int64) ([]*CourseItem, error)
	CountCourses(ctx context.Context) (int64, error)
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	ListPaymentsByStatus(ctx context.Context, status []int, limit, offset int64) ([]*Payment, error)
	CountPaymentsByStatus(ctx context.Context, status []int) (int64, error)
}
