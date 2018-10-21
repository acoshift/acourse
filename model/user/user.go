package user

// User type
type User struct {
	ID       string
	Role     Role
	Username string
	Name     string
	Email    string
	AboutMe  string
	Image    string
}

// Role type
type Role struct {
	Admin      bool
	Instructor bool
}

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

// IsEnroll checks is user enrolled a course
type IsEnroll struct {
	ID       string
	CourseID string

	Result bool
}

// Get gets user from id
type Get struct {
	ID string

	Result *User
}
