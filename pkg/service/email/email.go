package email

import (
	"log"

	"github.com/acoshift/acourse/pkg/acourse"
	_context "golang.org/x/net/context"
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
func (s *service) Send(ctx _context.Context, req *acourse.Email) (*acourse.Empty, error) {
	if len(req.To) == 0 {
		return new(acourse.Empty), nil
	}
	log.Printf("Send mail to %s\n", req.To)

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", req.To...)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	return nil, s.dialer.DialAndSend(m)
}
