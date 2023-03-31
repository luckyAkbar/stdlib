package mail_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/luckyAkbar/stdlib/helper"
	"github.com/luckyAkbar/stdlib/mail"
	mail_mock "github.com/luckyAkbar/stdlib/mail/mock"
	"github.com/stretchr/testify/assert"
)

func TestMailUtility_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)

	sendInBlue := mail_mock.NewMockClient(ctrl)
	mailgun := mail_mock.NewMockClient(ctrl)

	utility := mail.NewUtility(sendInBlue, mailgun)

	ctx := context.TODO()

	m := &mail.Mail{
		ID: helper.GenerateID(),
		To: []mail.GenericReceipient{
			{
				Name:  "test name",
				Email: "test.email@mail.test",
			},
		},
		Subject: "Testing",
	}

	sendInBlueMD := "metadatasendInBlue"
	mailgunMD := "metadatamailgun"

	t.Run("ok - first client", func(t *testing.T) {
		sendInBlue.EXPECT().SendEmail(ctx, m).Return(sendInBlueMD, nil)
		sendInBlue.EXPECT().GetClientName().Return(mail.SendInBlueSignature)

		metadata, signature, err := utility.SendEmail(ctx, m)

		assert.NoError(t, err)
		assert.Equal(t, metadata, sendInBlueMD)
		assert.Equal(t, mail.SendInBlueSignature, signature)
	})

	t.Run("ok - second client", func(t *testing.T) {
		sendInBlue.EXPECT().SendEmail(ctx, m).Return("", mail.ErrSendInBlueNotActivated)
		mailgun.EXPECT().SendEmail(ctx, m).Return(mailgunMD, nil)
		mailgun.EXPECT().GetClientName().Return(mail.MailgunSignature)

		metadata, signature, err := utility.SendEmail(ctx, m)

		assert.NoError(t, err)
		assert.Equal(t, metadata, mailgunMD)
		assert.Equal(t, mail.MailgunSignature, signature)
	})

	t.Run("all client err", func(t *testing.T) {
		sendInBlue.EXPECT().SendEmail(ctx, m).Return("", mail.ErrSendInBlueNotActivated)
		mailgun.EXPECT().SendEmail(ctx, m).Return("", mail.ErrMailgunNotActivated)

		_, _, err := utility.SendEmail(ctx, m)

		assert.Error(t, err)
	})
}
