package app

import "time"

// UserRawPayload type
type UserRawPayload struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Photo    *string `json:"photo"`
	AboutMe  *string `json:"aboutMe"`
}

// Validate validates UserRawPayload
func (x *UserRawPayload) Validate() error {
	return nil
}

// Payload builds UserPayload from UserRawPayload
func (x *UserRawPayload) Payload() *UserPayload {
	r := UserPayload{}
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

// CourseRawPayload type
type CourseRawPayload struct {
	Title            *string                    `json:"title"`
	ShortDescription *string                    `json:"shortDescription"`
	Description      *string                    `json:"description"`
	Photo            *string                    `json:"photo"`
	Start            *time.Time                 `json:"start"`
	Video            *string                    `json:"video"`
	Type             *string                    `json:"type"`
	Price            *float64                   `json:"price"`
	DiscountedPrice  *float64                   `json:"discountedPrice"`
	URL              *string                    `json:"url"`
	Contents         []*CourseContentRawPayload `json:"contents"`
	Enroll           *bool                      `json:"enroll"`
	Public           *bool                      `json:"public"`
	Attend           *bool                      `json:"attend"`
	Assignment       *bool                      `json:"assignment"`
	Purchase         *bool                      `json:"purchase"`
}

// Validate validates CourseRawPayload
func (x *CourseRawPayload) Validate() error {
	return nil
}

// Payload builds CoursePayload from CourseRawPayload
func (x *CourseRawPayload) Payload() *CoursePayload {
	r := CoursePayload{}
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
	if x.Type != nil {
		r.Type = *x.Type
	}
	r.Contents = make([]*CourseContentPayload, len(x.Contents))
	for i := range x.Contents {
		r.Contents[i] = x.Contents[i].Payload()
	}
	if x.Enroll != nil {
		r.Enroll = *x.Enroll
	}
	if x.Public != nil {
		r.Public = *x.Public
	}
	if x.Attend != nil {
		r.Attend = *x.Attend
	}
	if x.Assignment != nil {
		r.Assignment = *x.Assignment
	}
	if x.Purchase != nil {
		r.Purchase = *x.Purchase
	}
	return &r
}

// CourseContentRawPayload type
type CourseContentRawPayload struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Video       *string `json:"video"`
	DownloadURL *string `json:"downloadURL"`
}

// Validate validates CourseContentRawPayload
func (x *CourseContentRawPayload) Validate() error {
	return nil
}

// Payload builds CourseContentPayload from CourseContentRawPayload
func (x *CourseContentRawPayload) Payload() *CourseContentPayload {
	r := CourseContentPayload{}
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

// CourseEnrollRawPayload type
type CourseEnrollRawPayload struct {
	Code *string `json:"code"`
	URL  *string `json:"url"`
}

// Validate validates CourseEnrollRawPayload
func (x *CourseEnrollRawPayload) Validate() error {
	return nil
}

// Payload builds CourseEnrollPayload from CourseEnrollRawPayload
func (x *CourseEnrollRawPayload) Payload() *CourseEnrollPayload {
	r := CourseEnrollPayload{}
	if x.Code != nil {
		r.Code = *x.Code
	}
	if x.URL != nil {
		r.URL = *x.URL
	}
	return &r
}
