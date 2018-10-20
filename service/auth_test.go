package service_test

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/firebase"

	. "github.com/acoshift/acourse/service"
)

var _ = Describe("Auth", func() {
	var (
		s    Service
		repo *mockRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = &mockRepo{}
		s = New(Config{
			Repository: repo,
		})
	})

	Describe("SignUp", func() {
		It("should error with zero value email", func() {
			userID, err := s.SignUp(ctx, "", "123456")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with invalid email", func() {
			userID, err := s.SignUp(ctx, "invalid email", "123456")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with zero value password", func() {
			userID, err := s.SignUp(ctx, "test@test.com", "")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with too short password", func() {
			userID, err := s.SignUp(ctx, "test@test.com", "123")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with too long password", func() {
			userID, err := s.SignUp(ctx, "test@test.com", strings.Repeat("1", 100))

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should propagate error when firebase error", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				return fmt.Errorf("error")
			})

			userID, err := s.SignUp(ctx, "test@test.com", "12345678")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when email not available", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				m.Result = "123"
				return nil
			})
			repo.On("RegisterUser", &RegisterUser{
				ID:       "123",
				Username: "123",
				Email:    "test@test.com",
			}).Return(entity.ErrEmailNotAvailable)

			userID, err := s.SignUp(ctx, "test@test.com", "12345678")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when username not available", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				m.Result = "123"
				return nil
			})
			repo.On("RegisterUser", &RegisterUser{
				ID:       "123",
				Username: "123",
				Email:    "test@test.com",
			}).Return(entity.ErrUsernameNotAvailable)

			userID, err := s.SignUp(ctx, "test@test.com", "12345678")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when database return error", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				m.Result = "123"
				return nil
			})
			repo.On("RegisterUser", &RegisterUser{
				ID:       "123",
				Username: "123",
				Email:    "test@test.com",
			}).Return(fmt.Errorf("db error"))

			userID, err := s.SignUp(ctx, "test@test.com", "12345678")

			Expect(err).NotTo(BeNil())
			Expect(userID).To(BeZero())
		})

		It("should return user id when success", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				m.Result = "123"
				return nil
			})
			repo.On("RegisterUser", &RegisterUser{
				ID:       "123",
				Username: "123",
				Email:    "test@test.com",
			}).Return(nil)

			userID, err := s.SignUp(ctx, "test@test.com", "12345678")

			Expect(err).To(BeNil())
			Expect(userID).To(Equal("123"))
		})
	})

	Describe("SignInPassword", func() {
		It("should error when sign in with zero value email and password", func() {
			userID, err := s.SignInPassword(ctx, "", "")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with valid email but zero password", func() {
			userID, err := s.SignInPassword(ctx, "test@test.com", "")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with zero email but valid password", func() {
			userID, err := s.SignInPassword(ctx, "", "123456")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with valid email but wrong password", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.VerifyPassword) error {
				Expect(m.Email).To(Equal("test@test.com"))
				Expect(m.Password).To(Equal("fakepass"))
				return fmt.Errorf("invalid")
			})

			userID, err := s.SignInPassword(ctx, "test@test.com", "fakepass")

			Expect(err).NotTo(BeNil())
			Expect(userID).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.VerifyPassword) error {
				Expect(m.Email).To(Equal("test@test.com"))
				Expect(m.Password).To(Equal("123456"))
				m.Result = "aqswde"
				return nil
			})
			repo.On("IsUserExists", "aqswde").Return(true, nil)

			userID, err := s.SignInPassword(ctx, "test@test.com", "123456")

			Expect(err).To(BeNil())
			Expect(userID).To(Equal("aqswde"))
		})
	})

	Describe("SendPasswordResetEmail", func() {
		It("should error when email is empty", func() {
			err := s.SendPasswordResetEmail(ctx, "")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should success when email not found in firebase", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.GetUserByEmail) error {
				Expect(m.Email).To(Equal("notfound@test.com"))
				return fmt.Errorf("not found")
			})

			err := s.SendPasswordResetEmail(ctx, "notfound@test.com")

			Expect(err).To(BeNil())
		})

		It("should propagate error when firebase send email error", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.GetUserByEmail) error {
				Expect(m.Email).To(Equal("test@test.com"))
				m.Result = &admin.UserRecord{
					DisplayName:   "tester",
					Email:         "test@test.com",
					EmailVerified: true,
					UserID:        "12345",
				}
				return nil
			})
			dispatcher.Register(func(_ context.Context, m *firebase.SendPasswordResetEmail) error {
				Expect(m.Email).To(Equal("test@test.com"))
				return fmt.Errorf("some error")
			})

			err := s.SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should success when email valid", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.GetUserByEmail) error {
				Expect(m.Email).To(Equal("test@test.com"))
				m.Result = &admin.UserRecord{
					DisplayName:   "tester",
					Email:         "test@test.com",
					EmailVerified: true,
					UserID:        "12345",
				}
				return nil
			})
			dispatcher.Register(func(_ context.Context, m *firebase.SendPasswordResetEmail) error {
				Expect(m.Email).To(Equal("test@test.com"))
				return nil
			})

			err := s.SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).To(BeNil())
		})
	})
})
