package handler

import (
	"strings"

	"github.com/rkrohk/gobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Commandhandler handles bot commands
func Commandhandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.Message.Command() {
	case "help":
		{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "Hello, This is Rohan's bot"
			go bot.Send(msg)
		}

	case "echo":
		{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = strings.Replace(update.Message.Text, "/echo", "", 1)
			go bot.Send(msg)
		}

	case "slap":
		{
			currentUser := update.Message.From.FirstName
			replyToMessage := update.Message.ReplyToMessage
			victimUser := ""
			if replyToMessage != nil {
				victimUser = replyToMessage.From.FirstName
			} else {
				messageWithoutCommand := strings.Replace(update.Message.Text, "/slap ", "", 1)
				victimUser = messageWithoutCommand
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpers.Slap(currentUser, victimUser))
			msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
			go bot.Send(msg)
		}

	case "timer":
		go helpers.Timer(bot, &update)
	}
}
