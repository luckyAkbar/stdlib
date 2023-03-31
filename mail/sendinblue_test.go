package mail_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luckyAkbar/stdlib/helper"
	"github.com/luckyAkbar/stdlib/mail"
	"github.com/sendinblue/APIv3-go-library/lib"
	"github.com/stretchr/testify/assert"
)

func TestSendInBlue(t *testing.T) {
	m := &mail.Mail{
		To: []mail.GenericReceipient{
			{
				Name:  "test name",
				Email: "email23@gmail.test",
			},
		},
	}

	t.Run("ok - send email", func(t *testing.T) {
		mockSrc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}))

		sibCfg := &lib.Configuration{
			BasePath:      mockSrc.URL,
			DefaultHeader: make(map[string]string),
			UserAgent:     "Swagger-Codegen/1.0.0/go",
		}
		sibCfg.AddDefaultHeader("api-key", helper.GenerateID())

		sibClient := lib.NewAPIClient(sibCfg)

		sib := &mail.SendInBlue{}
		sib.Set(sibClient, &lib.SendSmtpEmailSender{}, true)
		_, err := sib.SendEmail(context.TODO(), m)

		assert.NoError(t, err)
	})

	t.Run("sib not activated", func(t *testing.T) {
		sib := &mail.SendInBlue{}
		sib.Set(nil, nil, false)
		_, err := sib.SendEmail(context.TODO(), m)

		assert.Error(t, err)
		assert.Equal(t, err, mail.ErrSendInBlueNotActivated)
	})

	t.Run("err", func(t *testing.T) {
		sibCfg := &lib.Configuration{
			BasePath:      "invalid url",
			DefaultHeader: make(map[string]string),
			UserAgent:     "Swagger-Codegen/1.0.0/go",
		}
		sibCfg.AddDefaultHeader("api-key", helper.GenerateID())

		sibClient := lib.NewAPIClient(sibCfg)

		sib := &mail.SendInBlue{}
		sib.Set(sibClient, &lib.SendSmtpEmailSender{}, true)
		_, err := sib.SendEmail(context.TODO(), m)

		assert.Error(t, err)
	})

	t.Run("err - received non 201", func(t *testing.T) {
		mockSrc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		sibCfg := &lib.Configuration{
			BasePath:      mockSrc.URL,
			DefaultHeader: make(map[string]string),
			UserAgent:     "Swagger-Codegen/1.0.0/go",
		}
		sibCfg.AddDefaultHeader("api-key", helper.GenerateID())

		sibClient := lib.NewAPIClient(sibCfg)

		sib := &mail.SendInBlue{}
		sib.Set(sibClient, &lib.SendSmtpEmailSender{}, true)
		_, err := sib.SendEmail(context.TODO(), m)

		assert.Error(t, err)
	})
}
