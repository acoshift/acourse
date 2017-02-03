package email

import (
	"context"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"gopkg.in/gomail.v2"
)

// New creates new email service server
func New(config Config) acourse.EmailServiceServer {
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

// Send sends an email
func (s *service) Send(ctx context.Context, req *acourse.Email) (*acourse.Empty, error) {
	if len(req.GetTo()) == 0 {
		return new(acourse.Empty), nil
	}
	internal.InfoLogger.Printf("Send mail to %s\n", req.GetTo())

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", req.GetTo()...)
	m.SetHeader("Subject", req.GetSubject())
	m.SetBody("text/html", req.GetBody())

	err := s.dialer.DialAndSend(m)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}
