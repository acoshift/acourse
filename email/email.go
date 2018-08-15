package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// Sender is the email sender
type Sender interface {
	Send(to, subject, body string) error
}

// NewSMTPSender creates new smtp sender
func NewSMTPSender(c SMTPConfig) Sender {
	return &smtpSender{
		Dialer: gomail.NewPlainDialer(c.Server, c.Port, c.User, c.Password),
		From:   c.From,
	}
}

// SMTPConfig type
type SMTPConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	From     string
}

type smtpSender struct {
	Dialer *gomail.Dialer
	From   string
}

func (s *smtpSender) Send(to, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("invalid to")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.Dialer.DialAndSend(m)
}
