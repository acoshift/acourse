package user

// Create creates new user
type Create struct {
	ID       string
	Username string
	Name     string
	Email    string
	Image    string
}

// Update updates user
type Update struct {
	ID       string
	Username string
	Name     string
	AboutMe  string
}

// IsExists checks is user exists
type IsExists struct {
	ID string

	Result bool
}

// SetImage sets user image
type SetImage struct {
	ID    string
	Image string
}
