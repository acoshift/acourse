package app

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func sendEmail(to string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("invalid to")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return emailDialer.DialAndSend(m)
}
