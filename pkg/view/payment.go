package view

import "time"

// Payment type
type Payment struct {
	ID            string      `json:"id"`
	User          *UserTiny   `json:"user"`
	Course        *CourseMini `json:"course"`
	OriginalPrice float64     `json:"originalPrice"`
	Price         float64     `json:"price"`
	Code          string      `json:"code"`
	URL           string      `json:"url"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	At            time.Time   `json:"at"`
}

// PaymentCollection type
type PaymentCollection []*Payment
