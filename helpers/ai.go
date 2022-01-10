package helpers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/helpers/ai"
)

//SendMessage sends a message to the bot to store
func SendMessage(message string) {

	go ai.SendMessage(message)

}

//GetMessage gets message from the bot and replies to the user
func GetMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message

	response, err := ai.GetMessageResponse(message)
	if err != nil {
		log.Println("Error getting response", err)
		return
	}
	replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, response.Message)
	replyMessage.ReplyToMessageID = message.MessageID
	bot.Send(replyMessage)
}
