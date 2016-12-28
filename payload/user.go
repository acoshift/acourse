package payload

// User type
type User struct {
	Username string
	Name     string
	Photo    string
	AboutMe  string
}

// RawUser type
type RawUser struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Photo    *string `json:"photo"`
	AboutMe  *string `json:"aboutMe"`
}

// Validate validates model
func (x *RawUser) Validate() error {
	return nil
}

// Payload builds UserPayload from model
func (x *RawUser) Payload() *User {
	r := User{}
	if x.Username != nil {
		r.Username = *x.Username
	}
	if x.Name != nil {
		r.Name = *x.Name
	}
	if x.Photo != nil {
		r.Photo = *x.Photo
	}
	if x.AboutMe != nil {
		r.AboutMe = *x.AboutMe
	}
	return &r
}
