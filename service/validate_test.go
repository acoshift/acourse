package service

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateImage(t *testing.T) {
	Convey("Given empty file header", t, func() {
		fh := &multipart.FileHeader{}

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then ui error should be return", func() {
				So(err, ShouldNotBeNil)
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})

	Convey("Given html file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "text/html")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then ui error should be return", func() {
				So(err, ShouldNotBeNil)
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})

	Convey("Given jpg file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/jpg")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then no error should be return", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given jpeg file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/jpeg")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then no error should be return", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given png file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/png")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then no error should be return", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given bmp file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/bmp")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then ui error should be return", func() {
				So(err, ShouldNotBeNil)
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})

	Convey("Given svg file header", t, func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/svg+xml")

		Convey("When validate", func() {
			err := validateImage(fh)

			Convey("Then ui error should be return", func() {
				So(err, ShouldNotBeNil)
				So(IsUIError(err), ShouldBeTrue)
			})
		})
	})
}
