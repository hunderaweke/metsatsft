package pkg

import (
	gomail "gopkg.in/mail.v2"

	"github.com/hunderaweke/metsasft/config"
)

func sendEmail(subject, text, html, sender string, recipients []string, config config.Config) error {
	message := gomail.NewMessage()
	message.SetHeader("From", sender)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)
	message.AddAlternative("text/html", html)
	dialer := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
