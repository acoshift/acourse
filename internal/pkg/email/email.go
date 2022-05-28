package email

import (
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/acoshift/acourse/internal/pkg/config"
)

func Init() {
	dialer = gomail.NewDialer(
		config.String("email_server"),
		config.Int("email_port"),
		config.String("email_user"),
		config.String("email_password"),
	)
	from = config.String("email_from")
}

var (
	dialer *gomail.Dialer
	from   string
)

// Send send an email
func Send(to, subject, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("invalid to")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return dialer.DialAndSend(msg)
}
