package app

import (
	"time"
)

// UserView type
type UserView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	AboutMe  string `json:"aboutMe"`
}

// UserTinyView type
type UserTinyView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
}

// UserMeView type
type UserMeView struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	AboutMe  string    `json:"aboutMe"`
	Role     *RoleView `json:"role"`
}

// RoleView type
type RoleView struct {
	Admin      bool `json:"admin"`
	Instructor bool `json:"instructor"`
}

// CourseView type
type CourseView struct {
}

// CoursePublicView type
type CoursePublicView struct {
}

// CourseTinyView type
type CourseTinyView struct {
	ID               string        `json:"id"`
	Owner            *UserTinyView `json:"owner"`
	Title            string        `json:"title"`
	ShortDescription string        `json:"shortDescription"`
	Photo            string        `json:"photo"`
	Start            time.Time     `json:"start"`
	URL              string        `json:"url"`
	Type             string        `json:"type"`
	Price            float64       `json:"price"`
	DiscountedPrice  *float64      `json:"discountedPrice,omitempty"`
	Student          int           `json:"student"`
}

// CourseTinyCollectionView type
type CourseTinyCollectionView []*CourseTinyView
