package payload

import (
	"time"
)

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

// RawCourse type
type RawCourse struct {
	Title            *string             `json:"title"`
	ShortDescription *string             `json:"shortDescription"`
	Description      *string             `json:"description"`
	Photo            *string             `json:"photo"`
	Start            *time.Time          `json:"start"`
	Video            *string             `json:"video"`
	Price            *float64            `json:"price"`
	DiscountedPrice  *float64            `json:"discountedPrice"`
	URL              *string             `json:"url"`
	Contents         []*RawCourseContent `json:"contents"`
	Attend           *bool               `json:"attend"`
	Assignment       *bool               `json:"assignment"`
}

// Validate validates model
func (x *RawCourse) Validate() error {
	return nil
}

// Payload builds CoursePayload from model
func (x *RawCourse) Payload() *Course {
	r := Course{}
	if x.Title != nil {
		r.Title = *x.Title
	}
	if x.ShortDescription != nil {
		r.ShortDescription = *x.ShortDescription
	}
	if x.Description != nil {
		r.Description = *x.Description
	}
	if x.Photo != nil {
		r.Photo = *x.Photo
	}
	if x.Start != nil {
		r.Start = *x.Start
	}
	if x.Video != nil {
		r.Video = *x.Video
	}
	r.Contents = make([]*CourseContent, len(x.Contents))
	for i := range x.Contents {
		r.Contents[i] = x.Contents[i].Payload()
	}
	if x.Attend != nil {
		r.Attend = *x.Attend
	}
	if x.Assignment != nil {
		r.Assignment = *x.Assignment
	}
	return &r
}

// RawCourseContent type
type RawCourseContent struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Video       *string `json:"video"`
	DownloadURL *string `json:"downloadURL"`
}

// Validate validates model
func (x *RawCourseContent) Validate() error {
	return nil
}

// Payload builds CourseContentPayload from model
func (x *RawCourseContent) Payload() *CourseContent {
	r := CourseContent{}
	if x.Title != nil {
		r.Title = *x.Title
	}
	if x.Description != nil {
		r.Description = *x.Description
	}
	if x.Video != nil {
		r.Video = *x.Video
	}
	if x.DownloadURL != nil {
		r.DownloadURL = *x.DownloadURL
	}
	return &r
}

// RawCourseEnroll type
type RawCourseEnroll struct {
	Code *string `json:"code"`
	URL  *string `json:"url"`
}

// Validate validates model
func (x *RawCourseEnroll) Validate() error {
	return nil
}

// Payload builds CourseEnrollPayload from model
func (x *RawCourseEnroll) Payload() *CourseEnroll {
	r := CourseEnroll{}
	if x.Code != nil {
		r.Code = *x.Code
	}
	if x.URL != nil {
		r.URL = *x.URL
	}
	return &r
}
