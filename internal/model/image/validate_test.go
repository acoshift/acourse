package image_test

import (
	"mime/multipart"
	"net/textproto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/acoshift/acourse/internal/model/app"

	. "github.com/acoshift/acourse/internal/model/image"
)

var _ = Describe("Validate", func() {
	It("should return error when file is empty", func() {
		fh := &multipart.FileHeader{}
		err := Validate(fh)

		Expect(err).ToNot(BeNil())
		Expect(app.IsUIError(err)).To(BeTrue())
	})

	It("should return error when file is html", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "text/html")
		err := Validate(fh)

		Expect(err).ToNot(BeNil())
		Expect(app.IsUIError(err)).To(BeTrue())
	})

	It("should success when file is jpg", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/jpg")
		err := Validate(fh)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should success when file is jpeg", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/jpeg")
		err := Validate(fh)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should success when file is png", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/png")
		err := Validate(fh)

		Expect(err).NotTo(HaveOccurred())
	})

	It("should return error when file is bmp", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/bmp")
		err := Validate(fh)

		Expect(err).ToNot(BeNil())
		Expect(app.IsUIError(err)).To(BeTrue())
	})

	It("should return error when file is svg", func() {
		fh := &multipart.FileHeader{}
		fh.Header = make(textproto.MIMEHeader)
		fh.Header.Set("Content-Type", "image/svg+xml")
		err := Validate(fh)

		Expect(err).ToNot(BeNil())
		Expect(app.IsUIError(err)).To(BeTrue())
	})
})
