package pkg

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	gomail "gopkg.in/mail.v2"

	"github.com/hunderaweke/metsasft/config"
)

func sendEmail(subject, html, sender string, recipients []string, config config.Config) error {
	message := gomail.NewMessage()
	message.SetHeader("From", sender)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", subject)
	message.AddAlternative("text/html", html)
	dialer := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func SendResetEmail(recipient, token string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	link := fmt.Sprintf(config.Server.Url+"/reset-password?token=%s&email=%s", token, recipient)
	file, err := os.Open("templates/reset_password.html")
	defer file.Close()
	if err != nil {
		return err
	}
	html, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	parsedHtml, err := template.New("reset_password").Parse(string(html))
	if err != nil {
		return err
	}
	var buf strings.Builder
	err = parsedHtml.Execute(&buf, struct {
		Link string
	}{
		Link: link,
	})
	if err != nil {
		return err
	}
	htmlText := buf.String()
	err = sendEmail("Reset Password", htmlText, "hunderaweke@gmail.com", []string{recipient}, config)
	if err != nil {
		return err
	}
	return nil
}
