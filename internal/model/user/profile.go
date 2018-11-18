package user

import "mime/multipart"

// UpdateProfile command
type UpdateProfile struct {
	ID       string
	Username string
	Name     string
	AboutMe  string
	Image    *multipart.FileHeader
}
