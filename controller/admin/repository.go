package admin

import (
	"context"

	"github.com/acoshift/acourse/entity"
)

// Repository is the admin storage
type Repository interface {
	ListUsers(ctx context.Context, limit, offset int64) ([]*entity.UserItem, error)
	CountUsers(ctx context.Context) (int64, error)
	ListCourses(ctx context.Context, limit, offset int64) ([]*entity.Course, error)
	CountCourses(ctx context.Context) (int64, error)
	GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error)
	ListPaymentsByStatus(ctx context.Context, statuses []int, limit, offset int64) ([]*entity.Payment, error)
	CountPaymentsByStatuses(ctx context.Context, statuses []int) (int64, error)
}
