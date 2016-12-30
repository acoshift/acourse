package ctrl

import (
	"acourse/model"
	"acourse/store"
	"fmt"
	"log"
	"time"
)

// StartNotiPayment runs notification payment service
// send email notification to admin when waiting payment existing in database
func StartNotiPayment(db *store.DB) {
	go func() {
		for {
			// check is payments have status waiting
			log.Println("Run Notification Payment")
			payments, err := db.PaymentList(model.PaymentStatusWaiting)
			if err == nil && len(payments) > 0 {
				err = SendMail(Email{
					To:      []string{"acoshift@gmail.com", "k.chalermsook@gmail.com"},
					Subject: "Acourse - Payment Received",
					Body:    fmt.Sprintf("%d payments pending", len(payments)),
				})
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(2 * time.Hour)
		}
	}()
}
