package acourse

import (
	"time"

	"github.com/acoshift/acourse/pkg/model"
)

// ToUser builds an User message from an User model
func ToUser(x *model.User) *User {
	return &User{
		Id:       x.ID,
		Username: x.Username,
		Name:     x.Name,
		Photo:    x.Photo,
		AboutMe:  x.AboutMe,
	}
}

// ToUsers builds repeated User message from User models
func ToUsers(xs model.Users) []*User {
	rs := make([]*User, len(xs))
	for i, x := range xs {
		rs[i] = ToUser(x)
	}
	return rs
}

// ToUserTiny builds an User tiny message from an User model
func ToUserTiny(x *model.User) *UserTiny {
	return &UserTiny{
		Id:       x.ID,
		Username: x.Username,
		Name:     x.Username,
		Photo:    x.Photo,
	}
}

// ToUsersTiny builds repeated User tiny message from User models
func ToUsersTiny(xs model.Users) []*UserTiny {
	rs := make([]*UserTiny, len(xs))
	for i, x := range xs {
		rs[i] = ToUserTiny(x)
	}
	return rs
}

// ToRole builds a Role message from Role model
func ToRole(x *model.Role) *Role {
	return &Role{
		Admin:      x.Admin,
		Instructor: x.Instructor,
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ToPayment builds a Payment message from Payment model
func ToPayment(x *model.Payment) *Payment {
	return &Payment{
		Id:            x.ID,
		CreatedAt:     formatTime(x.CreatedAt),
		UpdatedAt:     formatTime(x.UpdatedAt),
		UserId:        x.UserID,
		CourseId:      x.CourseID,
		OriginalPrice: x.OriginalPrice,
		Price:         x.Price,
		Code:          x.Code,
		Url:           x.URL,
		Status:        string(x.Status),
		At:            formatTime(x.At),
	}
}

// ToPayments builds repeated Payment message from Payment models
func ToPayments(xs model.Payments) []*Payment {
	rs := make([]*Payment, len(xs))
	for i, x := range xs {
		rs[i] = ToPayment(x)
	}
	return rs
}

// ToCourse builds a Course message from Course model
func ToCourse(x *model.Course) *Course {
	return &Course{}
}

// ToCourses builds repeated Course message from Course models
func ToCourses(xs model.Courses) []*Course {
	rs := make([]*Course, len(xs))
	for i, x := range xs {
		rs[i] = ToCourse(x)
	}
	return rs
}

// ToCourseTiny builds a Course tiny message from Course model
func ToCourseTiny(x *model.Course) *CourseTiny {
	return &CourseTiny{
		Id:    x.ID,
		Title: x.Title,
	}
}

// ToCoursesTiny builds repeated Course tiny message from Course models
func ToCoursesTiny(xs model.Courses) []*CourseTiny {
	rs := make([]*CourseTiny, len(xs))
	for i, x := range xs {
		rs[i] = ToCourseTiny(x)
	}
	return rs
}
