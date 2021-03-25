package helpers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/sed"
)

func Sed(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	message := update.Message

	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	reply.ParseMode = "markdown"

	repliedToMessage := update.Message.ReplyToMessage

	if repliedToMessage != nil {
		if replyString, err := sed.Sed(repliedToMessage.Text, message.Text); err != nil {
			reply.Text = "Error replacing text"
			log.Println("Sed.go error ", err)

		} else {
			reply.Text = "_Did you mean?_\n\n" + replyString
			reply.ReplyToMessageID = repliedToMessage.MessageID
		}

	} else {
		reply.Text = "Please reply to a message"
	}

	go bot.Send(reply)
}
