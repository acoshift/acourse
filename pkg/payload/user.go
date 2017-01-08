package payload

// User type
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// Validate validates model
func (x *User) Validate() error {
	return nil
}
