package view

import (
	"time"

	"github.com/acoshift/acourse/pkg/model"
)

// Course view
type Course struct {
	ID               string         `json:"id"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	Title            string         `json:"title"`
	ShortDescription string         `json:"shortDescription"`
	Description      string         `json:"description"`
	Photo            string         `json:"photo"`
	Owner            string         `json:"owner"`
	Start            time.Time      `json:"start"`
	URL              string         `json:"url"`
	Video            string         `json:"video"`
	Type             string         `json:"type"`
	Price            float64        `json:"price"`
	DiscountedPrice  float64        `json:"discountedPrice"`
	Contents         CourseContents `json:"courseContents"`
	EnrollDetail     string         `json:"enrollDetail"`
	Options          CourseOption   `json:"options"`
}

// ToCourse builds Course view from Course model
func ToCourse(x *model.Course) *Course {
	return &Course{
		x.ID,
		x.CreatedAt,
		x.UpdatedAt,
		x.Title,
		x.ShortDescription,
		x.Description,
		x.Photo,
		x.Owner,
		x.Start,
		x.URL,
		x.Video,
		string(x.Type),
		x.Price,
		x.DiscountedPrice,
		ToCourseContents(x.Contents),
		x.EnrollDetail,
		*ToCourseOption(&x.Options),
	}
}

// CourseOption view
type CourseOption struct {
	Public    bool `json:"public"`
	Enroll    bool `json:"enroll"`
	Attend    bool `json:"attend"`
	Assigment bool `json:"assignment"`
	Discount  bool `json:"discount"`
}

// ToCourseOption builds Course option view from Course option model
func ToCourseOption(x *model.CourseOption) *CourseOption {
	return &CourseOption{x.Public, x.Enroll, x.Attend, x.Assignment, x.Discount}
}

// CourseOptionPublic view
type CourseOptionPublic struct {
	Enroll   bool `json:"enroll"`
	Discount bool `json:"discount"`
}

// ToCourseOptionPublic builds Course option public view from Course option model
func ToCourseOptionPublic(x *model.CourseOption) *CourseOptionPublic {
	return &CourseOptionPublic{x.Enroll, x.Discount}
}

// CourseContent view
type CourseContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
	DownloadURL string `json:"downloadURL"`
}

// ToCourseContent builds Course content view from Course content model
func ToCourseContent(x *model.CourseContent) *CourseContent {
	return &CourseContent{x.Title, x.Description, x.Video, x.DownloadURL}
}

// CourseContents view
type CourseContents []*CourseContent

// ToCourseContents builds Course contents view from Course contents model
func ToCourseContents(xs model.CourseContents) CourseContents {
	rs := make(CourseContents, len(xs))
	for i, x := range xs {
		rs[i] = ToCourseContent(&x)
	}
	return rs
}
