package payload

import (
	"fmt"
	"time"
)

// Course type
type Course struct {
	Title            string           `json:"title"`
	ShortDescription string           `json:"shortDescription"`
	Description      string           `json:"description"`
	Photo            string           `json:"photo"`
	Start            time.Time        `json:"start"`
	Video            string           `json:"video"`
	Contents         []*CourseContent `json:"contents"`
	Attend           bool             `json:"attend"`
	Assignment       bool             `json:"assignment"`
}

// CourseContent type
type CourseContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Video       string `json:"video"`
	DownloadURL string `json:"downloadURL"`
}

// CourseEnroll type
type CourseEnroll struct {
	Code  string  `json:"code"`
	URL   string  `json:"url"`
	Price float64 `json:"price"`
}

// Validate validates model
func (x *Course) Validate() error {
	return nil
}

// Validate validates model
func (x *CourseEnroll) Validate() error {
	if x.Price < 0 {
		return fmt.Errorf("payload: price should be 0 or above")
	}
	return nil
}
