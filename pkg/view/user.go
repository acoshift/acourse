package view

// User type
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// UserTiny type
type UserTiny struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// UserMe type
type UserMe struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
	Role     *Role  `json:"role"`
}
