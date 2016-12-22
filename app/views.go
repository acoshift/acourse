package app

// UserView type
type UserView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// UserTinyView type
type UserTinyView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// UserMeView type
type UserMeView struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	AboutMe  string    `json:"aboutMe"`
	Role     *RoleView `json:"role"`
}

// RoleView type
type RoleView struct {
	Admin      bool `json:"admin"`
	Instructor bool `json:"instructor"`
}
