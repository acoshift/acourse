package appctx

import (
	"context"
	"testing"

	"github.com/acoshift/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppCtx(t *testing.T) {
	Convey("Given empty session context", t, func() {
		ctx := context.Background()
		ctx = newSessionContext(ctx, &session.Session{})

		Convey("Should get non-nil session", func() {
			So(getSession(ctx), ShouldNotBeNil)
		})

		Convey("Should get non-nil flash", func() {
			So(GetFlash(ctx), ShouldNotBeNil)
		})

		Convey("Should get empty flash", func() {
			So(GetFlash(ctx).Count(), ShouldEqual, 0)
		})

		Convey("Should get zero user id", func() {
			So(GetUserID(ctx), ShouldBeZeroValue)
		})

		Convey("Should get nil user", func() {
			So(GetUser(ctx), ShouldBeNil)
		})

		Convey("Should get zero open id session", func() {
			So(GetOpenIDState(ctx), ShouldBeZeroValue)
		})

		Convey("When set user id", func() {
			SetUserID(ctx, "1")

			Convey("Should be able to get that user id", func() {
				So(GetUserID(ctx), ShouldEqual, "1")
			})
		})

		Convey("When set open id state", func() {
			SetOpenIDState(ctx, "s1")

			Convey("Should be able to get that state", func() {
				So(GetOpenIDState(ctx), ShouldEqual, "s1")

				Convey("When delete that state", func() {
					DelOpenIDState(ctx)

					Convey("Should not able to get that state", func() {
						So(GetOpenIDState(ctx), ShouldBeZeroValue)
					})
				})
			})
		})
	})
}
