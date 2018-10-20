package service_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/service"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*service.User, error) {
	args := m.Mock.Called(email)
	return args.Get(0).(*service.User), args.Error(1)
}

func (m *mockRepo) RegisterUser(ctx context.Context, x *service.RegisterUser) error {
	args := m.Mock.Called(x)
	return args.Error(0)
}

func (m *mockRepo) UpdateUser(ctx context.Context, x *service.UpdateUser) error {
	args := m.Mock.Called(x)
	return args.Error(0)
}

func (m *mockRepo) SetUserImage(ctx context.Context, userID string, image string) error {
	args := m.Mock.Called(userID, image)
	return args.Error(0)
}

func (m *mockRepo) IsUserExists(ctx context.Context, userID string) (exists bool, err error) {
	args := m.Mock.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) RegisterCourse(ctx context.Context, x *service.RegisterCourse) (courseID string, err error) {
	args := m.Mock.Called(x)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
	args := m.Mock.Called(courseID)
	return args.Get(0).(*entity.Course), args.Error(1)
}

func (m *mockRepo) UpdateCourse(ctx context.Context, x *service.UpdateCourseModel) error {
	args := m.Mock.Called(x)
	return args.Error(0)
}

func (m *mockRepo) SetCourseImage(ctx context.Context, courseID string, image string) error {
	args := m.Mock.Called(courseID, image)
	return args.Error(0)
}

func (m *mockRepo) SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error {
	args := m.Mock.Called(courseID, x)
	return args.Error(0)
}

func (m *mockRepo) RegisterCourseContent(ctx context.Context, x *entity.RegisterCourseContent) (contentID string, err error) {
	args := m.Mock.Called(x)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetCourseContent(ctx context.Context, contentID string) (*course.Content, error) {
	args := m.Mock.Called(contentID)
	return args.Get(0).(*course.Content), args.Error(1)
}

func (m *mockRepo) ListCourseContents(ctx context.Context, courseID string) ([]*course.Content, error) {
	args := m.Mock.Called(courseID)
	return args.Get(0).([]*course.Content), args.Error(1)
}

func (m *mockRepo) UpdateCourseContent(ctx context.Context, contentID, title, desc, videoID string) error {
	args := m.Mock.Called(contentID, title, desc, videoID)
	return args.Error(0)
}

func (m *mockRepo) DeleteCourseContent(ctx context.Context, contentID string) error {
	args := m.Mock.Called(contentID)
	return args.Error(0)
}

func (m *mockRepo) RegisterPayment(ctx context.Context, x *service.RegisterPayment) error {
	args := m.Mock.Called(x)
	return args.Error(0)
}

func (m *mockRepo) GetPayment(ctx context.Context, paymentID string) (*service.Payment, error) {
	args := m.Mock.Called(paymentID)
	return args.Get(0).(*service.Payment), args.Error(1)
}

func (m *mockRepo) SetPaymentStatus(ctx context.Context, paymentID string, status int) error {
	args := m.Mock.Called(paymentID, status)
	return args.Error(0)
}

func (m *mockRepo) HasPendingPayment(ctx context.Context, userID string, courseID string) (bool, error) {
	args := m.Mock.Called(userID, courseID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) RegisterEnroll(ctx context.Context, userID string, courseID string) error {
	args := m.Mock.Called(userID, courseID)
	return args.Error(0)
}

func (m *mockRepo) IsEnrolled(ctx context.Context, userID string, courseID string) (bool, error) {
	args := m.Mock.Called(userID, courseID)
	return args.Bool(0), args.Error(1)
}
