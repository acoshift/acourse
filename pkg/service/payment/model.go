package payment

import (
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/ds"
)

const kindPayment = "Payment"

type payment struct {
	ds.StringIDModel
	ds.StampModel
	UserID        string
	CourseID      string
	OriginalPrice float64 `datastore:",noindex"`
	Price         float64 `datastore:",noindex"`
	Code          string
	URL           string `datastore:",noindex"`
	Status        status
	At            time.Time
}

type status string

const (
	statusWaiting  status = "waiting"
	statusApproved status = "approved"
	statusRejected status = "rejected"
)

// Approve approves a payment
func (x *payment) Approve() {
	x.Status = statusApproved
	x.At = time.Now()
}

// Reject rejects a payment
func (x *payment) Reject() {
	x.Status = statusRejected
	x.At = time.Now()
}

func toPayment(x *payment) *acourse.Payment {
	return &acourse.Payment{
		Id:            x.ID(),
		CreatedAt:     x.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     x.UpdatedAt.Format(time.RFC3339),
		UserId:        x.UserID,
		CourseId:      x.CourseID,
		OriginalPrice: x.OriginalPrice,
		Price:         x.Price,
		Code:          x.Code,
		Url:           x.URL,
		Status:        string(x.Status),
		At:            x.At.Format(time.RFC3339),
	}
}

func toPayments(xs []*payment) []*acourse.Payment {
	rs := make([]*acourse.Payment, len(xs))
	for i, x := range xs {
		rs[i] = toPayment(x)
	}
	return rs
}
