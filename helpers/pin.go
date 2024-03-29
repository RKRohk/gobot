package helpers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//Pin function pins a message
func Pin(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message

	silent := strings.Contains(message.Text, "silent")

	repliedToMessage := message.ReplyToMessage

	replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if repliedToMessage == nil {
		replyMessage.Text = "Please reply to a message"
		go bot.Send(replyMessage)
		return
	}

	pinMessageConfig := tgbotapi.PinChatMessageConfig{ChatID: message.Chat.ID, MessageID: repliedToMessage.MessageID, DisableNotification: silent}
	_, err := bot.Send(pinMessageConfig)
	if err != nil {
		replyMessage.Text = "There was an error pinning the message"
		go bot.Send(replyMessage)
		log.Println(err)
	}
}

func Unpin(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	repliedToMessage := update.Message.ReplyToMessage
	if repliedToMessage == nil {
		replyMessage.Text = "Please reply to a message"
	} else {
		unpinMessageConfig := tgbotapi.UnpinChatMessageConfig{ChatID: message.Chat.ID}
		go bot.Send(unpinMessageConfig)
	}
}
