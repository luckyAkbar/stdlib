package mail

import (
	"context"
	"fmt"
)

// ClientSignature signature for every registered mailing client
type ClientSignature string

// Utility mail utility interface
type Utility interface {
	// SendEmail send email using any available mailinng client. Returning metadata, client signature and error
	// will retry using the next available client if the previous returning error
	SendEmail(ctx context.Context, mail *Mail) (string, ClientSignature, error)
}

// Client must be implemented by any client to be registered in mail utility
type Client interface {
	// SendEmail send email, returning the metadata, and error if any
	SendEmail(ctx context.Context, mail *Mail) (string, error)

	// GetClientName client name for mail client
	GetClientName() ClientSignature
}

type mail struct {
	clients []Client
}

// NewUtility return new mail utility
func NewUtility(clients ...Client) Utility {
	return &mail{
		clients,
	}
}

func (m *mail) SendEmail(ctx context.Context, mail *Mail) (metadata string, sig ClientSignature, err error) {
	for _, client := range m.clients {
		metadata, err = client.SendEmail(ctx, mail)

		if err == nil {
			return metadata, client.GetClientName(), nil
		}
	}

	return "", "", fmt.Errorf("failed to send email: %w", err)
}
