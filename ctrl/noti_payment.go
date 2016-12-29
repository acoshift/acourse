package ctrl

import (
	"acourse/model"
	"acourse/store"
	"fmt"
	"time"
)

// StartNotiPayment runs notification payment service
// send email notification to admin when waiting payment existing in database
func StartNotiPayment(db *store.DB) {
	go func() {
		for {
			// check is payments have status waiting
			payments, err := db.PaymentList(model.PaymentStatusWaiting)
			if err == nil && len(payments) > 0 {
				SendMail(Email{
					To:      []string{"contact@acourse.io"},
					Subject: "Admin Notification",
					Body:    fmt.Sprintf("%d payments waiting for action", len(payments)),
				})
			}
			time.Sleep(2 * time.Hour)
		}
	}()
}
