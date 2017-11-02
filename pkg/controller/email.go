package controller

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func (c *ctrl) sendEmail(to string, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("invalid to")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", c.emailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	err := c.emailDialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
