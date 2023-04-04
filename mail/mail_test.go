package mail_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/helper"
	"github.com/sweet-go/stdlib/mail"
	"gopkg.in/guregu/null.v4"
)

func TestMail(t *testing.T) {
	m := mail.Mail{
		ID: helper.GenerateID(),
		To: []mail.GenericReceipient{
			{
				Name:  "test name",
				Email: "test email",
			},
		},
		Cc: []mail.GenericReceipient{
			{
				Name:  "test name",
				Email: "test email",
			},
		},
		Bcc: []mail.GenericReceipient{
			{
				Name:  "test name",
				Email: "test email",
			},
		},
		HTMLContent: "test html content",
		Subject:     "test subject",
		Metadata:    null.StringFrom("test metadata"),
	}

	t.Run("send in blue to", func(t *testing.T) {
		res := m.SendInBlueTo()

		assert.Equal(t, len(res), 1)
	})

	t.Run("send in blue cc", func(t *testing.T) {
		res := m.SendInBlueCc()

		assert.Equal(t, len(res), 1)
	})

	t.Run("send in blue bcc", func(t *testing.T) {
		res := m.SendInBlueBcc()

		assert.Equal(t, len(res), 1)
	})

	t.Run("mailgun to", func(t *testing.T) {
		res := m.MailgunTo()

		assert.Equal(t, len(res), 1)
	})

	t.Run("mailgun cc", func(t *testing.T) {
		res := m.MailgunCC()

		assert.Equal(t, len(res), 1)
	})

	t.Run("mailgun bcc", func(t *testing.T) {
		res := m.MailgunBCC()

		assert.Equal(t, len(res), 1)
	})
}
