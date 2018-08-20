package service_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/acoshift/go-firebase-admin"

	. "github.com/acoshift/acourse/service"
)

var _ = Describe("Auth", func() {
	var (
		s    Service
		repo *mockRepo
		auth *mockAuth
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = &mockRepo{}
		auth = &mockAuth{}
		s = New(Config{
			Repository: repo,
			Auth:       auth,
		})
	})

	Describe("SignUp", func() {
		It("should error when sign up with zero value email", func() {
			userID, err := s.SignUp(ctx, "", "123456")

			Expect(err).NotTo(BeNil())
			Expect(userID).To(BeZero())
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
			auth.On("VerifyPassword", "test@test.com", "fakepass").Return("", fmt.Errorf("invalid"))

			userID, err := s.SignInPassword(ctx, "test@test.com", "fakepass")

			Expect(err).NotTo(BeNil())
			Expect(userID).To(BeZero())
		})

		It("should success when sign in with valid email and password", func() {
			auth.On("VerifyPassword", "test@test.com", "123456").Return("aqswde", nil)
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
			auth.On("GetUserByEmail", "notfound@test.com").Return((*firebase.UserRecord)(nil), fmt.Errorf("not found"))

			err := s.SendPasswordResetEmail(ctx, "notfound@test.com")

			Expect(err).To(BeNil())
		})

		It("should propagate error when firebase send email error", func() {
			auth.On("GetUserByEmail", "test@test.com").Return(&firebase.UserRecord{
				DisplayName:   "tester",
				Email:         "test@test.com",
				EmailVerified: true,
				UserID:        "12345",
			}, nil)
			auth.On("SendPasswordResetEmail", "test@test.com").Return(fmt.Errorf("some error"))

			err := s.SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).NotTo(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should success when email valid", func() {
			auth.On("GetUserByEmail", "test@test.com").Return(&firebase.UserRecord{
				DisplayName:   "tester",
				Email:         "test@test.com",
				EmailVerified: true,
				UserID:        "12345",
			}, nil)
			auth.On("SendPasswordResetEmail", "test@test.com").Return(nil)

			err := s.SendPasswordResetEmail(ctx, "test@test.com")

			Expect(err).To(BeNil())
		})
	})
})
