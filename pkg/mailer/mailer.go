package mailer

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
)

type mailer struct {
	dialer *mail.Dialer
}

//go:generate mockgen -source=mailer.go -destination=mock/mailer.go -package=mock
type MailerInterface interface {
	DialAndSend(msg *mail.Message) error
}

func NewMailer() *mailer {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	// Initialize a new mail.Dialer instance with the given SMTP server settings. We
	// also configure this to use a 10-second timeout whenever we send an email.
	dialer := mail.NewDialer(os.Getenv("MAIL_HOST"), port, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"))
	dialer.Timeout = 10 * time.Second
	return &mailer{
		dialer: dialer,
	}
}

func (m *mailer) DialAndSend(msg *mail.Message) error {
	return m.dialer.DialAndSend(msg)
}
