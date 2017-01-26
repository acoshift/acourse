package acourse

import (
	"time"

	"github.com/acoshift/acourse/pkg/model"
)

// ToUserTiny builds an User tiny message from an User model
func ToUserTiny(x *User) *User {
	return &User{
		Id:       x.Id,
		Username: x.Username,
		Name:     x.Username,
		Photo:    x.Photo,
	}
}

// ToUsersTiny builds repeated User tiny message from User models
func ToUsersTiny(xs []*User) []*User {
	rs := make([]*User, len(xs))
	for i, x := range xs {
		rs[i] = ToUserTiny(x)
	}
	return rs
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ToCourse builds a Course message from Course model
func ToCourse(x *model.Course) *Course {
	r := &Course{
		Id:               x.ID(),
		CreatedAt:        formatTime(x.CreatedAt),
		UpdatedAt:        formatTime(x.UpdatedAt),
		Title:            x.Title,
		ShortDescription: x.ShortDescription,
		Description:      x.Description,
		Photo:            x.Photo,
		Owner:            x.Owner,
		Start:            formatTime(x.Start),
		Url:              x.URL,
		Type:             string(x.Type),
		Video:            x.Video,
		Price:            x.Price,
		DiscountedPrice:  x.DiscountedPrice,
		Options: &Course_Option{
			Public:     x.Options.Public,
			Enroll:     x.Options.Enroll,
			Attend:     x.Options.Attend,
			Assignment: x.Options.Assignment,
			Discount:   x.Options.Discount,
		},
		EnrollDetail: x.EnrollDetail,
	}

	r.Contents = make([]*Course_Content, len(x.Contents))
	for i, c := range x.Contents {
		r.Contents[i] = &Course_Content{
			Title:       c.Title,
			Description: c.Description,
			Video:       c.Video,
			DownloadURL: c.DownloadURL,
		}
	}

	return r
}

// ToCourses builds repeated Course message from Course models
func ToCourses(xs model.Courses) []*Course {
	rs := make([]*Course, len(xs))
	for i, x := range xs {
		rs[i] = ToCourse(x)
	}
	return rs
}

// ToCourseSmall builds a Course small message from a Course model
func ToCourseSmall(x *model.Course) *CourseSmall {
	return &CourseSmall{
		Id:               x.ID(),
		Title:            x.Title,
		ShortDescription: x.ShortDescription,
		Photo:            x.Photo,
		Owner:            x.Owner,
		Start:            formatTime(x.Start),
		Url:              x.URL,
		Type:             string(x.Type),
		Price:            x.Price,
		DiscountedPrice:  x.DiscountedPrice,
		Options: &CourseSmall_Option{
			Discount: x.Options.Discount,
		},
	}
}

// ToCoursesSmall builds repeated Course small message from Course models
func ToCoursesSmall(xs model.Courses) []*CourseSmall {
	rs := make([]*CourseSmall, len(xs))
	for i, x := range xs {
		rs[i] = ToCourseSmall(x)
	}
	return rs
}

// ToCourseTiny builds a Course tiny message from Course model
func ToCourseTiny(x *model.Course) *CourseTiny {
	return &CourseTiny{
		Id:    x.ID(),
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
