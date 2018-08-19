package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSignInPassword(t *testing.T) {
	ctx := context.Background()

	Convey("Given zero value email", t, func() {
		email := ""

		Convey("Given zero value password", func() {
			password := ""

			Convey("When call the service", func() {
				s := &svc{}
				userID, err := s.SignInPassword(ctx, email, password)

				Convey("Then ui error should be return", func() {
					So(err, ShouldNotBeNil)
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

				Convey("Then ui error should be return", func() {
					So(err, ShouldNotBeNil)
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

				Convey("Then ui error should be return", func() {
					So(err, ShouldNotBeNil)
					So(IsUIError(err), ShouldBeTrue)
				})

				Convey("Then user id should be zero value", func() {
					So(userID, ShouldBeZeroValue)
				})
			})
		})
	})
}
