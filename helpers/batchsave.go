package helpers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/helpers/middleware"
	"github.com/rkrohk/gobot/helpers/note"
)

type BatchSaveType struct {
	Data data
}

type data struct {
	Tag   string
	Notes []*note.Note
}

//BatchSave handler adds the command to the user session
func BatchSave(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	handler := &BatchSaveType{}
	middleware.AddSession(middleware.UserChatID(update.Message.Chat.ID), "batchsave", handler)
	handler.Start(bot, update)
}

//Start function checks whether tag is present or not
//and then sets the data in the session
func (b *BatchSaveType) Start(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	messageWithoutCommand := update.Message.CommandArguments()
	tag := ExtractTag(messageWithoutCommand)
	if len(tag) <= 3 {
		reply.Text = "Please give a valid tag"

	} else {
		b.Data.Tag = tag
		log.Println("added session")
		reply.Text = "Please send the files you want to save\nSend /done to finish and /cancel to cancel"
	}
	bot.Send(reply)
}

//Continue function adds the document to the session
func (b *BatchSaveType) Continue(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *middleware.UserSession) {
	message := update.Message
	if message.Document != nil {
		file := message.Document
		note := &note.Note{FileID: file.FileID, Content: file.FileName, Tag: b.Data.Tag}
		b.Data.Notes = append(b.Data.Notes, note)
	} else {
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Please send a file")
		bot.Send(reply)
	}
}

//Done saves the notes and then indexes them all
func (b *BatchSaveType) Done(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *middleware.UserSession) {
	notes := b.Data.Notes
	err := note.BulkSaveNotes(notes)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if err != nil {
		log.Println("Error inserting notes", err)
		reply.Text = "There was an error saving your notes"
	} else {
		note.BulkIndexNotes(bot, notes)
		reply.Text = fmt.Sprintf("Note saved!\nUse `/get %s` to retrieve this note", b.Data.Tag)
		reply.ParseMode = "markdown"
	}
	bot.Send(reply)
}
