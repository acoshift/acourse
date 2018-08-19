package editor

import (
	"context"

	"github.com/acoshift/acourse/entity"
)

// Repository is editor storage
type Repository interface {
	GetCourse(ctx context.Context, courseID string) (*entity.Course, error)
	GetCourseUserID(ctx context.Context, courseID string) (string, error)
	GetCourseURL(ctx context.Context, courseID string) (string, error)
}
