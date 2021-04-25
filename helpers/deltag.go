package helpers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/note"
)

func DelTag(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message.CommandArguments()
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tag := ExtractTag(message)

	count, err := note.DeleteTag(tag)

	//TODO() also need to delete index from elasticsearch

	if err != nil || count == 0 {
		reply.Text = "This tag doesn't exist -__-"
	} else {
		reply.Text = fmt.Sprintf("Okay %d notes of tag `%s` have been deleted", count, tag)
		reply.ParseMode = "markdown"
	}

	bot.Send(reply)

}
