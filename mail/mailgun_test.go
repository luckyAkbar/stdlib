package mail_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luckyAkbar/stdlib/mail"
	mailgun "github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
)

func TestMailgun(t *testing.T) {
	t.Run("mg is not activated", func(t *testing.T) {
		mg := mail.Mailgun{}
		mg.Set(nil, "", false)

		_, err := mg.SendEmail(context.TODO(), nil)

		assert.Error(t, err)
		assert.Equal(t, err, mail.ErrMailgunNotActivated)
	})

	t.Run("err mailgun", func(t *testing.T) {
		mockSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		mgClient := mailgun.MailgunImpl{}
		mgClient.SetAPIBase(mockSrv.URL)

		mg := mail.Mailgun{}
		mg.Set(&mgClient, "", true)

		_, err := mg.SendEmail(context.TODO(), &mail.Mail{
			To: []mail.GenericReceipient{
				{
					Name:  "test name",
					Email: "test email",
				},
			},
			Subject: "test subj",
		})

		assert.Error(t, err)
	})

	// TODO: unit test for success case sending email using mailgun
}
