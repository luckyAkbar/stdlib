// Package mail is the functionality to interact with external service to send email
package mail

import (
	"context"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

// MailgunSignature signature of mailgun client
const MailgunSignature ClientSignature = "mailgun client"

// MailgunConfig configuration for mailgun client
type MailgunConfig struct {
	Domain            string
	PrivateKey        string
	IsActivated       bool
	ServerSenderEmail string
}

// Mailgun :nodoc:
type Mailgun struct {
	client            *mailgun.MailgunImpl
	isActivated       bool
	serverSenderEmail string
}

// NewMailgunClient create new mailgun client
func NewMailgunClient(config MailgunConfig) *Mailgun {
	client := mailgun.NewMailgun(config.Domain, config.PrivateKey)
	return &Mailgun{
		client,
		config.IsActivated,
		config.ServerSenderEmail,
	}
}

// SendEmail send email using sendinblue
func (mg *Mailgun) SendEmail(ctx context.Context, mail *Mail) (string, error) {
	if !mg.isActivated {
		return "", ErrMailgunNotActivated
	}

	message := mg.client.NewMessage(mg.serverSenderEmail, mail.Subject, "", mail.MailgunTo()...)
	message.SetHtml(mail.HTMLContent)

	for _, email := range mail.MailgunCC() {
		message.AddCC(email)
	}

	for _, email := range mail.MailgunBCC() {
		message.AddBCC(email)
	}

	_, id, err := mg.client.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetClientName returning client name signature
func (mg *Mailgun) GetClientName() ClientSignature {
	return MailgunSignature
}

// Set is used for testing purpose
func (mg *Mailgun) Set(client *mailgun.MailgunImpl, serverSenderEmail string, isActivated bool) {
	mg.client = client
	mg.isActivated = isActivated
	mg.serverSenderEmail = serverSenderEmail
}
