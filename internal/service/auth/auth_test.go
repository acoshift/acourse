package auth_test

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/auth"
	"github.com/acoshift/acourse/internal/pkg/model/user"

	. "github.com/acoshift/acourse/internal/service/auth"
)

var _ = Describe("Auth", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		Init()
	})

	Describe("SignUp", func() {
		It("should error with zero value email", func() {
			q := auth.SignUp{
				Email:    "",
				Password: "123456",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with invalid email", func() {
			q := auth.SignUp{
				Email:    "invalid email",
				Password: "123456",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with zero value password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with too short password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "123",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error with too long password", func() {
			q := auth.SignUp{
				Email:    "test@test.com",
				Password: strings.Repeat("1", 100),
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should propagate error when firebase error", func() {
			firAuth.error = fmt.Errorf("error")

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when email not available", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return user.ErrEmailNotAvailable
			})

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when username not available", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return user.ErrUsernameNotAvailable
			})

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when database return error", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return fmt.Errorf("db error")
			})

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(q.Result).To(BeZero())
		})

		It("should return user id when success", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return nil
			})

			q := auth.SignUp{
				Email:    "test@test.com",
				Password: "12345678",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).NotTo(HaveOccurred())
			Expect(q.Result).To(Equal("123"))
		})
	})

	Describe("SignInPassword", func() {
		It("should error when sign in with zero value email and password", func() {
			q := auth.SignInPassword{}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with valid email but zero password", func() {
			q := auth.SignInPassword{
				Email: "test@test.com",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with zero email but valid password", func() {
			q := auth.SignInPassword{
				Password: "123456",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(q.Result).To(BeZero())
		})

		It("should error when sign in with valid email but wrong password", func() {
			firAuth.error = fmt.Errorf("invalid")

			q := auth.SignInPassword{
				Email:    "test@test.com",
				Password: "fakepass",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).To(HaveOccurred())
			Expect(q.Result).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.IsExists) error {
				m.Result = true
				return nil
			})

			q := auth.SignInPassword{
				Email:    "test@test.com",
				Password: "123456",
			}
			err := bus.Dispatch(ctx, &q)

			Expect(err).NotTo(HaveOccurred())
			Expect(q.Result).NotTo(BeEmpty())
		})
	})

	Describe("SendPasswordResetEmail", func() {
		It("should error when email is empty", func() {
			err := bus.Dispatch(ctx, &auth.SendPasswordResetEmail{})

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
		})

		It("should success when email not found in firebase", func() {
			firAuth.error = fmt.Errorf("not found")

			err := bus.Dispatch(ctx, &auth.SendPasswordResetEmail{
				Email: "notfound@test.com",
			})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should success when email valid", func() {
			firAuth.error = nil

			err := bus.Dispatch(ctx, &auth.SendPasswordResetEmail{
				Email: "test@test.com",
			})

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
