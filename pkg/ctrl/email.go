package ctrl

import (
	"log"

	"gopkg.in/gomail.v2"
)

// Email type
type Email struct {
	To      []string
	Subject string
	Body    string
}

// EmailConfig is the config for send email
type EmailConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	From     string
}

var (
	emailConfig = EmailConfig{}
	emailDialer *gomail.Dialer
)

// InitMail inits email config and dialer
func InitMail(config EmailConfig) {
	emailConfig = config
	emailDialer = gomail.NewDialer(emailConfig.Server, emailConfig.Port, emailConfig.User, emailConfig.Password)
}

// SendMail sends email
func SendMail(context Email) error {
	if len(context.To) == 0 {
		return nil
	}
	log.Printf("Send mail to %s\n", context.To)

	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.From)
	m.SetHeader("To", context.To...)
	m.SetHeader("Subject", context.Subject)
	m.SetBody("text/html", context.Body)

	return emailDialer.DialAndSend(m)
}
