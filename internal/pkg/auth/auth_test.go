package auth_test

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/acoshift/acourse/internal/pkg/app"
	. "github.com/acoshift/acourse/internal/pkg/auth"
)

var _ = Describe("Auth", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("SignUp", func() {
		It("should error with empty email", func() {
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

		It("should error with empty password", func() {
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

		It("should error when email not available", func() {
			userID, err := SignUp(ctx, "notavailable@test.com", "12345678")

			Expect(err).To(HaveOccurred())
			Expect(app.IsUIError(err)).To(BeTrue())
			Expect(userID).To(BeZero())
		})

		It("should return user id when success", func() {
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
			userID, err := SignInPassword(ctx, "test@test.com", "fakepass")

			Expect(err).To(HaveOccurred())
			Expect(userID).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
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
			err := SendPasswordResetEmail(ctx, "notfound@test.com")

			Expect(err).NotTo(HaveOccurred())
		})

		It("should success when email valid", func() {
			err := SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
