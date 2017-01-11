package email

import (
	"context"
	"log"

	"github.com/acoshift/acourse/pkg/app"
	"gopkg.in/gomail.v2"
)

// New creates new email service
func New(config Config) app.EmailService {
	return &service{
		config: config,
		dialer: gomail.NewDialer(
			config.Server,
			config.Port,
			config.User,
			config.Password,
		),
	}
}

type service struct {
	config Config
	dialer *gomail.Dialer
}

// Config is the config for send email
type Config struct {
	Server   string
	Port     int
	User     string
	Password string
	From     string
}

// SendEmail sends an email
func (s *service) SendEmail(ctx context.Context, req *app.EmailRequest) error {
	if len(req.To) == 0 {
		return nil
	}
	log.Printf("Send mail to %s\n", req.To)

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", req.To...)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	return s.dialer.DialAndSend(m)
}
