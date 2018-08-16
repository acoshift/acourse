package service

import (
	"context"
	"testing"

	"github.com/acoshift/acourse/entity"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSendSignInMagicLinkEmail(t *testing.T) {
	ctx := context.Background()

	Convey("Given zero value email", t, func() {
		email := ""

		Convey("When call the service", func() {
			s := &svc{}
			err := s.SendSignInMagicLinkEmail(ctx, email)

			Convey("Then error should be return", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error should be ui error", func() {
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})

	Convey("Given invalid email", t, func() {
		email := "invalid email yay!!"

		Convey("When call the service", func() {
			s := &svc{}
			err := s.SendSignInMagicLinkEmail(ctx, email)

			Convey("Then error should be return", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error should be ui error", func() {
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})
}

func TestSignInPassword(t *testing.T) {
	ctx := context.Background()

	Convey("Given zero value email", t, func() {
		email := ""

		Convey("Given zero value password", func() {
			password := ""

			Convey("When call the service", func() {
				s := &svc{}
				userID, err := s.SignInPassword(ctx, email, password)

				Convey("Then error should be return", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("Then error should be ui error", func() {
					So(IsUIError(err), ShouldBeTrue)
				})

				Convey("Then user id should be zero value", func() {
					So(userID, ShouldBeZeroValue)
				})
			})
		})

		Convey("Given valid password", func() {
			password := "hello 1234"

			Convey("When call the service", func() {
				s := &svc{}
				userID, err := s.SignInPassword(ctx, email, password)

				Convey("Then error should be return", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("Then error should be ui error", func() {
					So(IsUIError(err), ShouldBeTrue)
				})

				Convey("Then user id should be zero value", func() {
					So(userID, ShouldBeZeroValue)
				})
			})
		})
	})

	Convey("Given valid email", t, func() {
		email := "test@test.com"

		Convey("Given zero value password", func() {
			password := ""

			Convey("When call the service", func() {
				s := &svc{}
				userID, err := s.SignInPassword(ctx, email, password)

				Convey("Then error should be return", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("Then error should be ui error", func() {
					So(IsUIError(err), ShouldBeTrue)
				})

				Convey("Then user id should be zero value", func() {
					So(userID, ShouldBeZeroValue)
				})
			})
		})
	})
}

func TestSignInMagicLink(t *testing.T) {
	ctx := context.Background()

	Convey("Given zero value magic link", t, func() {
		link := ""

		Convey("When call the service", func() {
			s := &svc{}
			userID, err := s.SignInMagicLink(ctx, link)

			Convey("Then error should be return", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error should be ui error", func() {
				So(IsUIError(err), ShouldBeTrue)
			})

			Convey("Then user id should be zero value", func() {
				So(userID, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given not found link", t, func() {
		link := "some-not-found-link"
		repo := &mockRepo{}
		repo.On("FindMagicLink", ctx, link).Return("", entity.ErrNotFound)

		Convey("When call the service", func() {
			s := &svc{Config{Repository: repo}}
			userID, err := s.SignInMagicLink(ctx, link)

			Convey("Then error should be return", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error should be ui error", func() {
				So(IsUIError(err), ShouldBeTrue)
			})

			Convey("Then error should not be repository error", func() {
				So(err, ShouldNotEqual, entity.ErrNotFound)
			})

			Convey("Then user id should be zero value", func() {
				So(userID, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given valid link", t, func() {
		link := "valid-link"
		repo := &mockRepo{}
		repo.On("FindMagicLink", ctx, link).Return("123", nil)

		Convey("When call the service", func() {
			s := &svc{Config{Repository: repo}}
			userID, err := s.SignInMagicLink(ctx, link)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then user id should be valid", func() {
				So(userID, ShouldEqual, "123")
			})
		})
	})
}
