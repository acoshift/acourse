package service_test

import (
	"mime/multipart"
	"net/textproto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/acoshift/acourse/service"
)

var _ = Describe("Validate", func() {
	Describe("ValidateImage", func() {
		It("should return error when file is empty", func() {
			fh := &multipart.FileHeader{}
			err := ValidateImage(fh)

			Expect(err).ToNot(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should return error when file is html", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "text/html")
			err := ValidateImage(fh)

			Expect(err).ToNot(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should success when file is jpg", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "image/jpg")
			err := ValidateImage(fh)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should success when file is jpeg", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "image/jpeg")
			err := ValidateImage(fh)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should success when file is png", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "image/png")
			err := ValidateImage(fh)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when file is bmp", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "image/bmp")
			err := ValidateImage(fh)

			Expect(err).ToNot(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})

		It("should return error when file is svg", func() {
			fh := &multipart.FileHeader{}
			fh.Header = make(textproto.MIMEHeader)
			fh.Header.Set("Content-Type", "image/svg+xml")
			err := ValidateImage(fh)

			Expect(err).ToNot(BeNil())
			Expect(IsUIError(err)).To(BeTrue())
		})
	})
})
