package mailer

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"os"

	"github.com/go-mail/mail/v2"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
)

//go:embed "templates"
var templateFS embed.FS

type mailService struct {
	sender string
	mailer MailerInterface
}

func NewMailService() *mailService {
	mailer := NewMailer()

	// Return a Mailer instance containing the dialer and sender information.
	return &mailService{
		sender: os.Getenv("MAIL_SENDER"),
		mailer: mailer,
	}
}

//go:generate mockgen -source=mail_service.go -destination=mock/mail_service.go -package=mock
type MailService interface {
	Send(ctx context.Context, recipient, templateFile string, data interface{}) error
}

// Send function takes the recipient email addresses,
// the name of the file containing the templates, and any
// dynamic data for the templates as an interface{} parameter.
func (m *mailService) Send(ctx context.Context, recipient, templateFile string, data interface{}) error {
	// Use the ParseFS() method to parse the required template file from the embedded
	// file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	//setup message format
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Call the DialAndSend() method on the dialer, passing in the message to send. This
	// opens a connection to the SMTP server, sends the message, then closes the
	// connection. If there is a timeout, it will return a "dial tcp: i/o timeout"
	// error.
	err = m.mailer.DialAndSend(msg)
	if err != nil {
		logger.Logger(ctx).Error("error sending email : ", err.Error())
		return err
	}
	logger.Logger(ctx).Info("Sent email notification to user:", recipient)

	return nil
}
