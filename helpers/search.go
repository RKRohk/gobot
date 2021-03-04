package helpers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/search"
)

func Search(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	messageWithoutCommand := strings.Replace(update.Message.Text, "/search ", "", 1)
	splitString := strings.SplitN(messageWithoutCommand, " ", 2)
	log.Println(splitString)
	hashTag := splitString[0]
	query := splitString[1]

	hashTag = strings.Replace(hashTag, "#", "", 1)

	fileIDs := search.Search(hashTag, query)

	for _, id := range fileIDs {
		fileSendConfig := tgbotapi.NewDocumentShare(update.Message.Chat.ID, id)
		bot.Send(fileSendConfig)
	}

}
