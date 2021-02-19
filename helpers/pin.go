package helpers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	_, err := bot.PinChatMessage(pinMessageConfig)
	if err != nil {
		replyMessage.Text = "There was an error pinning the message"
		go bot.Send(replyMessage)
		log.Println(err)
	}
}
