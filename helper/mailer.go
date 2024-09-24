package helper

import (
	"log"
	"net/smtp"
	"skripsi/config"
	"strings"
)

type Sent interface {
	SendEmail(to, HTMLbody string) error
}

type Mailer struct {
	config config.SMTPConfig
}

func NewMailer(config config.SMTPConfig) Sent {
	return &Mailer{
		config: config,
	}
}

func (m *Mailer) SendEmail(to, HTMLbody string) error {
	from := m.config.SMTPUSER
	pass := m.config.SMTPPASS
	host := m.config.SMTPHOST
	port := m.config.SMTPPORT
	address := host + ":" + port

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, pass, host)

	// Construct the email message.
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + "subject" + "\r\n\r\n" +
		HTMLbody)

	// Send the email.
	err := smtp.SendMail(address, auth, from, strings.Split(to, ","), message)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully to", to)
	return nil
}
