package view

// User view
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}
