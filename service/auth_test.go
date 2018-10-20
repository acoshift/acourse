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
	"github.com/acoshift/acourse/model/auth"
	"github.com/acoshift/acourse/model/firebase"

	. "github.com/acoshift/acourse/service"
)

var _ = Describe("Auth", func() {
	var (
		repo *mockRepo
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = &mockRepo{}
		Init(Config{
			Repository: repo,
		})
	})

	Describe("SignUp", func() {
		It("should error with zero value email", func() {
			q := auth.SignUp{
				Email:    "",
				Password: "123456",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with invalid email", func() {
			q := auth.SignUp{
				Email:    "invalid email",
				Password: "123456",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with zero value password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with too short password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "123",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with too long password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: strings.Repeat("1", 100),
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should propagate error when firebase error", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.CreateUser) error {
				Expect(m.User.Email).To(Equal("test@test.com"))
				Expect(m.User.Password).To(Equal("12345678"))
				return fmt.Errorf("error")
			})

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
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

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
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

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
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

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(q.Result).To(BeZero())
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

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).NotTo(HaveOccurred())
			Expect(q.Result).To(Equal("123"))
		})
	})

	Describe("SignInPassword", func() {
		It("should error when sign in with zero value email and password", func() {
			q := auth.SignInPassword{}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with valid email but zero password", func() {
			q := auth.SignInPassword{
				Email: "test@test.com",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with zero email but valid password", func() {
			q := auth.SignInPassword{
				Password: "123456",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with valid email but wrong password", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.VerifyPassword) error {
				Expect(m.Email).To(Equal("test@test.com"))
				Expect(m.Password).To(Equal("fakepass"))
				return fmt.Errorf("invalid")
			})

			q := auth.SignInPassword{
				Email:    "test@test.com",
				Password: "fakepass",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(q.Result).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.VerifyPassword) error {
				Expect(m.Email).To(Equal("test@test.com"))
				Expect(m.Password).To(Equal("123456"))
				m.Result = "aqswde"
				return nil
			})
			repo.On("IsUserExists", "aqswde").Return(true, nil)

			q := auth.SignInPassword{
				Email:    "test@test.com",
				Password: "123456",
			}
			err := dispatcher.Dispatch(ctx, &q)

			Expect(err).NotTo(HaveOccurred())
			Expect(q.Result).To(Equal("aqswde"))
		})
	})

	Describe("SendPasswordResetEmail", func() {
		It("should error when email is empty", func() {
			err := dispatcher.Dispatch(ctx, &auth.SendPasswordResetEmail{})

			Expect(err).To(HaveOccurred())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should success when email not found in firebase", func() {
			dispatcher.Register(func(_ context.Context, m *firebase.GetUserByEmail) error {
				Expect(m.Email).To(Equal("notfound@test.com"))
				return fmt.Errorf("not found")
			})

			err := dispatcher.Dispatch(ctx, &auth.SendPasswordResetEmail{
				Email: "notfound@test.com",
			})

			Expect(err).NotTo(HaveOccurred())
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

			err := dispatcher.Dispatch(ctx, &auth.SendPasswordResetEmail{
				Email: "test@test.com",
			})

			Expect(err).To(HaveOccurred())
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

			err := dispatcher.Dispatch(ctx, &auth.SendPasswordResetEmail{
				Email: "test@test.com",
			})

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
