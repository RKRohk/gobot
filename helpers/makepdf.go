package helpers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func MakePDF(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	response := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	repliedToMessage := update.Message.ReplyToMessage
	if repliedToMessage != nil {
		startMessageID := repliedToMessage.MessageID
		endMessageID := update.Message.MessageID - 1

		for i := startMessageID; i <= endMessageID; i++ {
			
		}

	} else {
		response.Text = "Please reply to a message"
	}
}
