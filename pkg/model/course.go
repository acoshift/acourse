package model

import (
	"time"
)

// Course model
type Course struct {
	Base
	Stampable
	Title            string `datastore:",noindex"`
	ShortDescription string `datastore:",noindex"`
	Description      string `datastore:",noindex"` // Markdown
	Photo            string `datastore:",noindex"` // URL
	Owner            string
	Start            time.Time
	URL              string
	Type             CourseType
	Video            string `datastore:",noindex"` // Cover Video
	Price            float64
	DiscountedPrice  float64
	Options          CourseOption
	Contents         []CourseContent `datastore:",noindex"`
	EnrollDetail     string          `datastore:",noindex"`
}

// Courses model
type Courses []*Course

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool `datastore:",noindex"`
	Attend     bool `datastore:",noindex"`
	Assignment bool `datastore:",noindex"`
	Discount   bool
}

// CourseContent type
type CourseContent struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"` // Markdown
	Video       string `datastore:",noindex"` // Youtube ID
	DownloadURL string `datastore:",noindex"` // Video download link
}

// CourseType type
type CourseType string

// CourseType
const (
	CourseTypeLive  CourseType = "live"
	CourseTypeVideo CourseType = "video"
	CourseTypeEbook CourseType = "ebook"
)

// CourseView type
type CourseView int

// CourseView
const (
	CourseViewDefault CourseView = iota
	CourseViewPublic
	CourseViewTiny
	CourseViewMini
)

// SetView sets view to model
func (x *Course) SetView(v CourseView) {
	x.view = v
}

// SetView sets view to model
func (xs Courses) SetView(v CourseView) {
	for _, x := range xs {
		x.SetView(v)
	}
}

// Expose exposes model to view
func (x *Course) Expose() interface{} {
	if x.view == nil {
		return nil
	}
	switch x.view.(CourseView) {
	case CourseViewDefault:
		return map[string]interface{}{
			"id":               x.ID,
			"createdAt":        x.CreatedAt,
			"updatedAt":        x.UpdatedAt,
			"title":            x.Title,
			"shortDescription": x.ShortDescription,
			"description":      x.Description,
			"photo":            x.Photo,
			"owner":            x.Owner,
			"start":            x.Start,
			"url":              x.URL,
			"video":            x.Video,
			"type":             string(x.Type),
			"price":            x.Price,
			"discountedPrice":  x.DiscountedPrice,
			"contents": func() interface{} {
				rs := make([]interface{}, len(x.Contents))
				for i, x := range x.Contents {
					rs[i] = map[string]interface{}{
						"title":       x.Title,
						"description": x.Description,
						"video":       x.Video,
						"downloadURL": x.DownloadURL,
					}
				}
				return rs
			}(),
			"enrollDetail": x.EnrollDetail,
			"options": map[string]bool{
				"enroll":     x.Options.Enroll,
				"public":     x.Options.Public,
				"attend":     x.Options.Attend,
				"assignment": x.Options.Assignment,
				"discount":   x.Options.Discount,
			},
		}
	case CourseViewPublic:
		return map[string]interface{}{
			"id":               x.ID,
			"createdAt":        x.CreatedAt,
			"updatedAt":        x.UpdatedAt,
			"title":            x.Title,
			"shortDescription": x.ShortDescription,
			"description":      x.Description,
			"photo":            x.Photo,
			"owner":            x.Owner,
			"start":            x.Start,
			"url":              x.URL,
			"type":             string(x.Type),
			"price":            x.Price,
			"discountedPrice":  x.DiscountedPrice,
			"enrollDetail":     x.EnrollDetail,
			"options": map[string]bool{
				"enroll":   x.Options.Enroll,
				"discount": x.Options.Discount,
			},
		}
	case CourseViewTiny:
		return map[string]interface{}{
			"id":               x.ID,
			"title":            x.Title,
			"shortDescription": x.ShortDescription,
			"photo":            x.Photo,
			"owner":            x.Owner,
			"start":            x.Start,
			"url":              x.URL,
			"type":             string(x.Type),
			"price":            x.Price,
			"discountedPrice":  x.DiscountedPrice,
			"discount":         x.Options.Discount,
		}
	case CourseViewMini:
		return map[string]interface{}{
			"id":    x.ID,
			"title": x.Title,
		}
	default:
		return nil
	}
}

// Expose exposes model to view
func (xs Courses) Expose() interface{} {
	rs := make([]interface{}, len(xs))
	for i, x := range xs {
		rs[i] = x.Expose()
	}
	return rs
}

// ExposeMap exposes model as map
func (xs Courses) ExposeMap() interface{} {
	rs := map[string]interface{}{}
	for _, x := range xs {
		rs[x.ID] = x.Expose()
	}
	return rs
}
