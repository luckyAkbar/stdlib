package mail

import (
	"context"

	"github.com/resendlabs/resend-go"
)

const ResendSignature ClientSignature = "resend"

// api key = re_6iSWfe8R_2dTh7juMqTtWPtwp2JabBr1R

type ResendConfig struct {
	ApiKey            string
	ServerSenderEmail string
}

type ResendClient struct {
	client            *resend.Client
	serverSenderEmail string
}

func NewResendClient(config ResendConfig) *ResendClient {
	client := resend.NewClient(config.ApiKey)
	return &ResendClient{
		client:            client,
		serverSenderEmail: config.ServerSenderEmail,
	}
}

func (r *ResendClient) GetClientName() ClientSignature {
	return ResendSignature
}

func (r *ResendClient) SendEmail(ctx context.Context, mail *Mail) (string, error) {
	resp, err := r.client.Emails.Send(&resend.SendEmailRequest{
		From:    r.serverSenderEmail,
		To:      mail.ResendTo(),
		Subject: mail.Subject,
		Cc:      mail.ResendCc(),
		Bcc:     mail.ResendBcc(),
		Html:    mail.HTMLContent,
	})

	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
