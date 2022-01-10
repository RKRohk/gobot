package helpers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/helpers/note"
	"github.com/rkrohk/gobot/helpers/search"
)

func DelTag(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message.CommandArguments()
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	tag := ExtractTag(message)

	count, err := note.DeleteTag(tag)

	if err != nil || count == 0 {
		reply.Text = "This tag doesn't exist -__-"
	} else {
		search.RemoveIndex(tag)
		reply.Text = fmt.Sprintf("Okay %d notes of tag `%s` have been deleted", count, tag)
		reply.ParseMode = "markdown"
	}

	bot.Send(reply)

}
