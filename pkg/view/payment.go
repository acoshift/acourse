package view

import (
	"time"

	"github.com/acoshift/acourse/pkg/model"
)

// Payment type
type Payment struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	CourseID      string    `json:"courseId"`
	OriginalPrice float64   `json:"originalPrice"`
	Price         float64   `json:"price"`
	Code          string    `json:"code"`
	URL           string    `json:"url"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	At            time.Time `json:"at"`
}

// PaymentCollection type
type PaymentCollection []*Payment

// ToPayment builds a payment view from a payment model
func ToPayment(m *model.Payment) *Payment {
	return &Payment{
		ID:            m.ID,
		UserID:        m.UserID,
		CourseID:      m.CourseID,
		OriginalPrice: m.OriginalPrice,
		Price:         m.Price,
		Code:          m.Code,
		URL:           m.URL,
		Status:        string(m.Status),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		At:            m.At,
	}
}

// ToPaymentCollection builds a payment collection view from payment models
func ToPaymentCollection(ms []*model.Payment) PaymentCollection {
	rs := make(PaymentCollection, len(ms))
	for i := range ms {
		rs[i] = ToPayment(ms[i])
	}
	return rs
}
