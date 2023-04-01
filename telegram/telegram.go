// Package telegram will act mainly as wrapper of "github.com/PaulSonOfLars/gotgbot/v2"
package telegram

import (
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/sirupsen/logrus"
)

var defaultUpdaterOpts = &ext.UpdaterOpts{
	ErrorLog: log.New(os.Stdout, "telegram_bot: ", log.LUTC),
	DispatcherOpts: ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logrus.Error("experiencing error from dispatcher error handler: err", err)
			_, err = ctx.EffectiveMessage.Reply(
				b,
				"Oops, bot experiencing error. Please retry again later.",
				&gotgbot.SendMessageOpts{
					ReplyToMessageId: ctx.Message.MessageId,
				},
			)

			if err != nil {
				logrus.Error("failed to send error message from dispatcher error handler: ", err)
			}

			return ext.DispatcherActionNoop
		},
	},
}

// NewUpdater return updater with default configuration
func NewUpdater() ext.Updater {
	return NewUpdaterWithConfig(defaultUpdaterOpts)
}

// NewUpdaterWithConfig return updater with custom configuration
func NewUpdaterWithConfig(cfg *ext.UpdaterOpts) ext.Updater {
	return ext.NewUpdater(cfg)
}
