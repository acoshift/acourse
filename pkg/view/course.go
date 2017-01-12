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

// CourseCollection type
type CourseCollection []*Course

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

// ToCourse builds a course view from a course model
func ToCourse(m *model.Course) *Course {
	return &Course{
		ID:               m.ID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Description:      m.Description,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Video:            m.Video,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		Contents:         ToCourseContentCollection(m.Contents),
		EnrollDetail:     m.EnrollDetail,
		Enroll:           m.Options.Enroll,
		Public:           m.Options.Public,
		Attend:           m.Options.Attend,
		Assignment:       m.Options.Assignment,
		Discount:         m.Options.Discount,
	}
}

// ToCoursePublic builds a course view from a course model
func ToCoursePublic(m *model.Course) *Course {
	return &Course{
		ID:               m.ID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Description:      m.Description,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Video:            m.Video,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		EnrollDetail:     m.EnrollDetail,
		Enroll:           m.Options.Enroll,
		Discount:         m.Options.Discount,
	}
}

// ToCourseTiny builds a course tiny view from a course model
func ToCourseTiny(m *model.Course) *Course {
	return &Course{
		ID:               m.ID,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Photo:            m.Photo,
		Start:            m.Start,
		URL:              m.URL,
		Type:             string(m.Type),
		Price:            m.Price,
		DiscountedPrice:  m.DiscountedPrice,
		Discount:         m.Options.Discount,
	}
}

// ToCourseCollection builds a course collection from course models
func ToCourseCollection(ms []*model.Course) CourseCollection {
	rs := make(CourseCollection, len(ms))
	for i := range ms {
		rs[i] = ToCourse(ms[i])
	}
	return rs
}

func ToCourses(ms []*model.Course, f func(*model.Course) *Course) CourseCollection {
	rs := make(CourseCollection, len(ms))
	for i := range ms {
		rs[i] = f(ms[i])
	}
	return rs
}

// ToCourseContent builds a course content view from a course content model
func ToCourseContent(m *model.CourseContent) *CourseContent {
	return &CourseContent{
		Title:       m.Title,
		Description: m.Description,
		Video:       m.Video,
		DownloadURL: m.DownloadURL,
	}
}

// ToCourseContentCollection builds a course content collection view from course content models
func ToCourseContentCollection(ms []model.CourseContent) CourseContentCollection {
	rs := make(CourseContentCollection, len(ms))
	for i := range ms {
		rs[i] = ToCourseContent(&ms[i])
	}
	return rs
}

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
