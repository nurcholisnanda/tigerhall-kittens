package mailer

import "github.com/go-mail/mail/v2"

type mailer struct {
}

//go:generate mockgen -source=mailer.go -destination=mock/mailer.go -package=mock
type MailerInterface interface {
	DialAndSend(dialer *mail.Dialer, msg *mail.Message) error
}

func New() *mailer {
	return &mailer{}
}

func (m *mailer) DialAndSend(dialer *mail.Dialer, msg *mail.Message) error {
	return dialer.DialAndSend(msg)
}
