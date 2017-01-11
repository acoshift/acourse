package view

import (
	"time"

	"github.com/acoshift/acourse/pkg/model"
)

// Course type
type Course struct {
	ID               string                  `json:"id"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	Owner            *UserTiny               `json:"owner"`
	Title            string                  `json:"title"`
	ShortDescription string                  `json:"shortDescription"`
	Description      string                  `json:"description"`
	Photo            string                  `json:"photo"`
	Start            time.Time               `json:"start"`
	URL              string                  `json:"url"`
	Video            string                  `json:"video"`
	Type             string                  `json:"type"`
	Price            float64                 `json:"price"`
	DiscountedPrice  float64                 `json:"discountedPrice"`
	Student          int                     `json:"student"`
	Contents         CourseContentCollection `json:"contents"`
	EnrollDetail     string                  `json:"enrollDetail"`
	Enrolled         bool                    `json:"enrolled"`
	Enroll           bool                    `json:"enroll"`
	Public           bool                    `json:"public"`
	Owned            bool                    `json:"owned"`
	Attend           bool                    `json:"attend"`
	Assignment       bool                    `json:"assignment"`
	Discount         bool                    `json:"discount"`
}

// CoursePublic type
type CoursePublic struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Owner            *UserTiny `json:"owner"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"shortDescription"`
	Description      string    `json:"description"`
	Photo            string    `json:"photo"`
	Start            time.Time `json:"start"`
	URL              string    `json:"url"`
	Video            string    `json:"video"`
	Type             string    `json:"type"`
	Price            float64   `json:"price"`
	DiscountedPrice  float64   `json:"discountedPrice"`
	EnrollDetail     string    `json:"enrollDetail"`
	Student          int       `json:"student"`
	Enroll           bool      `json:"enroll"`
	Discount         bool      `json:"discount"`
	PurchaseStatus   string    `json:"purchaseStatus"`
}

// CourseTiny type
type CourseTiny struct {
	ID               string    `json:"id"`
	Owner            *UserTiny `json:"owner"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"shortDescription"`
	Photo            string    `json:"photo"`
	Start            time.Time `json:"start"`
	URL              string    `json:"url"`
	Type             string    `json:"type"`
	Price            float64   `json:"price"`
	DiscountedPrice  float64   `json:"discountedPrice"`
	Student          int       `json:"student"`
	Discount         bool      `json:"discount"`
}

// CourseMini type
type CourseMini struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// CourseMiniCollection type
type CourseMiniCollection []*CourseMini

// CourseTinyCollection type
type CourseTinyCollection []*CourseTiny

// CourseContent type
type CourseContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
	DownloadURL string `json:"downloadURL"`
}

// CourseContentCollection type
type CourseContentCollection []*CourseContent

// ToCourseMini builds a course mini view from a course model
func ToCourseMini(m *model.Course) *CourseMini {
	return &CourseMini{
		ID:    m.ID,
		Title: m.Title,
	}
}

// ToCourseMiniCollection builds a course mini collection view from course models
func ToCourseMiniCollection(ms []*model.Course) CourseMiniCollection {
	rs := make(CourseMiniCollection, len(ms))
	for i := range ms {
		rs[i] = ToCourseMini(ms[i])
	}
	return rs
}
