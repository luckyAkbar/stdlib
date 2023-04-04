package mail

import (
	"context"
	"fmt"
	"net/http"

	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
	"github.com/sweet-go/stdlib/helper"
)

// SendInBlueSignature send in blue client signature
const SendInBlueSignature ClientSignature = "sendinblue"

// SendInBlue send in blue client
type SendInBlue struct {
	client      *sendinblue.APIClient
	sender      *sendinblue.SendSmtpEmailSender
	isActivated bool
}

// NewSendInBlueClient creates a new SendInBlue client
func NewSendInBlueClient(sender *sendinblue.SendSmtpEmailSender, sibAPIKey string, isActivated bool) *SendInBlue {
	sibConfig := sendinblue.NewConfiguration()
	sibConfig.AddDefaultHeader("api-key", sibAPIKey)

	return &SendInBlue{
		client:      sendinblue.NewAPIClient(sibConfig),
		sender:      sender,
		isActivated: isActivated,
	}
}

// SendEmail sends an email. error if status code from sendinblue server is not 201
func (s *SendInBlue) SendEmail(ctx context.Context, mail *Mail) (string, error) {
	if !s.isActivated {
		return "", ErrSendInBlueNotActivated
	}

	body := sendinblue.SendSmtpEmail{
		Sender:      s.sender,
		To:          mail.SendInBlueTo(),
		Cc:          mail.SendInBlueCc(),
		Bcc:         mail.SendInBlueBcc(),
		HtmlContent: mail.HTMLContent,
		Subject:     mail.Subject,
	}

	email, res, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		return "", err
	}

	defer helper.WrapCloser(res.Body.Close)

	if res.StatusCode != http.StatusAccepted {
		e := fmt.Errorf("failed to send email send in blue: %v", res)
		return "", e
	}

	return helper.Dump(email), nil
}

// GetClientName return client name signature sendinblue
func (s *SendInBlue) GetClientName() ClientSignature {
	return SendInBlueSignature
}

// Set is used only for testing purposes to be able to set dependencies
func (s *SendInBlue) Set(client *sendinblue.APIClient, sender *sendinblue.SendSmtpEmailSender, isActivated bool) {
	s.client = client
	s.sender = sender
	s.isActivated = isActivated
}
