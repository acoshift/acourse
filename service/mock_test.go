package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/acoshift/acourse/entity"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	args := m.Mock.Called(ctx, linkID, userID)
	return args.Error(0)
}

func (m *mockRepo) FindMagicLink(ctx context.Context, linkID string) (string, error) {
	args := m.Mock.Called(ctx, linkID)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) CanAcquireMagicLink(ctx context.Context, email string) (bool, error) {
	args := m.Mock.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Mock.Called(ctx, email)
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockRepo) RegisterUser(ctx context.Context, x *RegisterUser) error {
	args := m.Mock.Called(ctx, x)
	return args.Error(0)
}

func (m *mockRepo) UpdateUser(ctx context.Context, x *UpdateUser) error {
	args := m.Mock.Called(ctx, x)
	return args.Error(0)
}

func (m *mockRepo) SetUserImage(ctx context.Context, userID string, image string) error {
	args := m.Mock.Called(ctx, userID, image)
	return args.Error(0)
}

func (m *mockRepo) IsUserExists(ctx context.Context, userID string) (exists bool, err error) {
	args := m.Mock.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) RegisterCourse(ctx context.Context, x *RegisterCourse) (courseID string, err error) {
	args := m.Mock.Called(ctx, x)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
	args := m.Mock.Called(ctx, courseID)
	return args.Get(0).(*entity.Course), args.Error(1)
}

func (m *mockRepo) UpdateCourse(ctx context.Context, x *UpdateCourseModel) error {
	args := m.Mock.Called(ctx, x)
	return args.Error(0)
}

func (m *mockRepo) SetCourseImage(ctx context.Context, courseID string, image string) error {
	args := m.Mock.Called(ctx, courseID, image)
	return args.Error(0)
}

func (m *mockRepo) SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error {
	args := m.Mock.Called(ctx, courseID, x)
	return args.Error(0)
}

func (m *mockRepo) RegisterCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error) {
	args := m.Mock.Called(ctx, x)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetCourseContent(ctx context.Context, contentID string) (*entity.CourseContent, error) {
	args := m.Mock.Called(ctx, contentID)
	return args.Get(0).(*entity.CourseContent), args.Error(1)
}

func (m *mockRepo) ListCourseContents(ctx context.Context, courseID string) ([]*entity.CourseContent, error) {
	args := m.Mock.Called(ctx, courseID)
	return args.Get(0).([]*entity.CourseContent), args.Error(1)
}

func (m *mockRepo) UpdateCourseContent(ctx context.Context, contentID, title, desc, videoID string) error {
	args := m.Mock.Called(ctx, contentID, title, desc, videoID)
	return args.Error(0)
}

func (m *mockRepo) DeleteCourseContent(ctx context.Context, contentID string) error {
	args := m.Mock.Called(ctx, contentID)
	return args.Error(0)
}

func (m *mockRepo) RegisterPayment(ctx context.Context, x *RegisterPayment) error {
	args := m.Mock.Called(ctx, x)
	return args.Error(0)
}

func (m *mockRepo) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	args := m.Mock.Called(ctx, paymentID)
	return args.Get(0).(*Payment), args.Error(1)
}

func (m *mockRepo) SetPaymentStatus(ctx context.Context, paymentID string, status int) error {
	args := m.Mock.Called(ctx, paymentID, status)
	return args.Error(0)
}

func (m *mockRepo) HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error) {
	args := m.Mock.Called(ctx, userID, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) RegisterEnroll(ctx context.Context, userID string, courseID string) error {
	args := m.Mock.Called(ctx, userID, courseID)
	return args.Error(0)
}

func (m *mockRepo) IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error) {
	args := m.Mock.Called(ctx, userID, courseID)
	return args.Bool(0), args.Error(1)
}
