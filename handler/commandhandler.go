package handler

import (
	"fmt"
	"strings"

	"github.com/rkrohk/gobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Commandhandler handles bot commands
func Commandhandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.Message.Command() {
	case "id":
		{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("`%d`", update.Message.Chat.ID))
			msg.ParseMode = "markdown"
			go bot.Send(msg)
		}
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
			// msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = strings.Replace(update.Message.Text, "/echo", "", 1)
			go bot.Send(msg)
		}

	case "slap":
		{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			currentUser := update.Message.From.FirstName
			replyToMessage := update.Message.ReplyToMessage
			victimUser := ""
			if replyToMessage != nil {
				victimUser = replyToMessage.From.FirstName
				msg.ReplyToMessageID = replyToMessage.MessageID
			} else {
				messageWithoutCommand := strings.Replace(update.Message.Text, "/slap ", "", 1)
				victimUser = messageWithoutCommand
			}
			msg.Text = helpers.Slap(currentUser, victimUser)
			go bot.Send(msg)
		}

	case "timer":
		go helpers.Timer(bot, &update)

	case "remind":
		go helpers.Remind(bot, &update)

	case "addslap":
		{
			message := update.Message.Text

			//Add slap string to databse
			slapString := strings.Replace(message, "/addslap ", "", 1)
			go helpers.AddSlapToDB(slapString)

			//Reply to the user. Does not handle errors yet
			//TODO() add error message
			replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Added new slap response")
			go bot.Send(replyMessage)
		}
	case "showslaps":
		{
			go helpers.ShowSlapStrings(bot, &update)
		}
	case "addslapsticker":
		{
			go helpers.AddSlapSticker(bot, &update)
		}
	case "stickerslap":
		{
			go helpers.SlapWithSticker(bot, &update)
		}

	case "save":
		{
			go helpers.SaveNote(bot, &update)
		}
	case "get":
		{
			go helpers.GetNotes(bot, &update)
		}
	case "delnote":
		{
			go helpers.DeleteNote(bot, &update)
		}

	case "notes":
		{
			go helpers.GetAllTags(bot, &update)
		}

	case "pin":
		{
			go helpers.Pin(bot, &update)
		}

	case "unpin":
		{
			go helpers.Unpin(bot, &update)
		}

	case "search":
		{
			go helpers.Search(bot, &update)
		}

	case "dadjoke":
		{
			go helpers.DadJoke(bot, &update)
		}

	case "sed":
		{
			go helpers.Sed(bot, &update)
		}
	}

}
