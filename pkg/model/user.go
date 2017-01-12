package model

// User model
type User struct {
	Base
	Stampable
	Username string
	Name     string `datastore:",noindex"`
	Photo    string `datastore:",noindex"`
	AboutMe  string `datastore:",noindex"`
}

// UserView type
type UserView int

// UserView
const (
	_ UserView = iota
	UserViewDefault
	UserViewTiny
)

// Users type
type Users []*User

// SetView sets view to model
func (x *User) SetView(v UserView) {
	x.view = v
}

// SetView sets view to model
func (xs Users) SetView(v UserView) {
	for _, x := range xs {
		x.SetView(v)
	}
}

// Expose exposes model
func (x *User) Expose() interface{} {
	if x.view == nil {
		return nil
	}
	switch x.view.(UserView) {
	case UserViewDefault:
		return map[string]interface{}{
			"id":       x.ID,
			"username": x.Username,
			"name":     x.Name,
			"photo":    x.Photo,
			"aboutMe":  x.AboutMe,
		}
	case UserViewTiny:
		return map[string]interface{}{
			"id":       x.ID,
			"username": x.Username,
			"name":     x.Name,
			"photo":    x.Photo,
		}
	default:
		return nil
	}
}

// Expose exposes model
func (xs Users) Expose() interface{} {
	rs := make([]interface{}, len(xs))
	for i, x := range xs {
		rs[i] = x.Expose()
	}
	return rs
}
