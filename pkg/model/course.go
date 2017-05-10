package model

// Course model
type Course struct {
	id string
}

// ID returns course id
func (x *Course) ID() string {
	return x.id
}
