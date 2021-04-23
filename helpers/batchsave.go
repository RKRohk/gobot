package helpers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/middleware"
)

type BatchSaveType struct {
	Data data
}

type data struct {
	Tag   string
	Notes []*Note
}

func BatchSave(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	handler := &BatchSaveType{}
	middleware.AddSession(middleware.UserChatID(update.Message.Chat.ID), "batchsave", handler)
	handler.Start(bot, update)
}

func (b *BatchSaveType) Start(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	messageWithoutCommand := update.Message.CommandArguments()
	tag := ExtractTag(messageWithoutCommand)
	b.Data.Tag = tag
	log.Println("added session")
}

func (b *BatchSaveType) Continue(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *middleware.UserSession) {
	message := update.Message
	if message.Document != nil {
		file := message.Document
		note := &Note{FileID: file.FileID, Content: file.FileName, Tag: b.Data.Tag}
		b.Data.Notes = append(b.Data.Notes, note)
	}
}

func (b *BatchSaveType) Done(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *middleware.UserSession) {

}
