package email

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/model/email"
)

// SMTPConfig type
type SMTPConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	From     string
}

// InitSMTP inits email using SMTP strategy
func InitSMTP(c SMTPConfig) {
	s := smtpSender{
		Dialer: gomail.NewPlainDialer(c.Server, c.Port, c.User, c.Password),
		From:   c.From,
	}

	bus.Register(s.Send)
}

type smtpSender struct {
	Dialer *gomail.Dialer
	From   string
}

func (s *smtpSender) Send(_ context.Context, m *email.Send) error {
	if len(m.To) == 0 {
		return fmt.Errorf("invalid to")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Body)

	return s.Dialer.DialAndSend(msg)
}
