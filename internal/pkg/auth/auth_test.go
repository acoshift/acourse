package auth_test

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/user"

	. "github.com/acoshift/acourse/internal/pkg/auth"
)

var _ = Describe("Auth", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("SignUp", func() {
		It("should error with zero value email", func() {
			userID, err := SignUp(ctx, "", "123456")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with invalid email", func() {
			userID, err := SignUp(ctx, "invalid email", "123456")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with zero value password", func() {
			userID, err := SignUp(ctx, "test@test.com", "")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with too short password", func() {
			userID, err := SignUp(ctx, "test@test.com", "123")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error with too long password", func() {
			userID, err := SignUp(ctx, "test@test.com", strings.Repeat("1", 100))

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should propagate error when firebase error", func() {
			firAuth.error = fmt.Errorf("error")

			userID, err := SignUp(ctx, "test@test.com", "12345678")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when email not available", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return user.ErrEmailNotAvailable
			})

			userID, err := SignUp(ctx, "test@test.com", "12345678")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when username not available", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return user.ErrUsernameNotAvailable
			})

			userID, err := SignUp(ctx, "test@test.com", "12345678")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when database return error", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return fmt.Errorf("db error")
			})

			userID, err := SignUp(ctx, "test@test.com", "12345678")

			Expect(err).To(HaveOccurred())
			Expect(userID).To(BeZero())
		})

		It("should return user id when success", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.Create) error {
				return nil
			})

			userID, err := SignUp(ctx, "test@test.com", "12345678")

			Expect(err).NotTo(HaveOccurred())
			Expect(userID).To(Equal("123"))
		})
	})

	Describe("SignInPassword", func() {
		It("should error when sign in with zero value email and password", func() {
			userID, err := SignInPassword(ctx, "", "")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with valid email but zero password", func() {
			userID, err := SignInPassword(ctx, "test@test.com", "")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with zero email but valid password", func() {
			userID, err := SignInPassword(ctx, "", "123456")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should error when sign in with valid email but wrong password", func() {
			firAuth.error = fmt.Errorf("invalid")

			userID, err := SignInPassword(ctx, "test@test.com", "fakepass")

			Expect(err).To(HaveOccurred())
			Expect(userID).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
			firAuth.error = nil

			bus.Register(func(_ context.Context, m *user.IsExists) error {
				m.Result = true
				return nil
			})

			userID, err := SignInPassword(ctx, "test@test.com", "123456")

			Expect(err).NotTo(HaveOccurred())
			Expect(userID).NotTo(BeEmpty())
		})
	})

	Describe("SendPasswordResetEmail", func() {
		It("should error when email is empty", func() {
			err := SendPasswordResetEmail(ctx, "")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
		})

		It("should success when email not found in firebase", func() {
			firAuth.error = fmt.Errorf("not found")

			err := SendPasswordResetEmail(ctx, "notfound@test.com")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should success when email valid", func() {
			firAuth.error = nil

			err := SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
