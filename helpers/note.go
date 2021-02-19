package helpers

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Note denotes how a note is saved in the database
type Note struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	FileID            string             `bson:"fileID,omitempty"`
	Tag               string             `bson:"tag,omitempty"`
	MessageID         int                `bson:"messageID,omitempty"`
	MessageFromChatID int64              `bson:"messageFromChatID,omitempty"`
	Content           string             `bson:"content,omitempty"`
}

//SaveNote saves a note
func SaveNote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	repliedToDocument := update.Message.ReplyToMessage.Document
	repliedToMessage := update.Message.ReplyToMessage
	if repliedToDocument == nil {
		//Handle null document
		if repliedToMessage == nil {
			return
		}
	}

	message := update.Message

	reg := regexp.MustCompile("#\\w+")
	tag := reg.FindString(message.Text)

	log.Println(tag)
	if tag == "" {
		//Handle this
		return
	}

	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	notesCollection := client.Database("bot").Collection("notes")
	log.Println(repliedToMessage)
	var newNote Note
	if repliedToDocument != nil {
		newNote = Note{FileID: repliedToDocument.FileID, Tag: tag, Content: repliedToDocument.FileName}

	} else {
		newNote = Note{Tag: tag, MessageID: repliedToMessage.MessageID, MessageFromChatID: repliedToMessage.Chat.ID, Content: repliedToMessage.Text}
	}
	log.Println(newNote)
	_, err = notesCollection.InsertOne(ctx, &newNote)
	replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if err != nil {
		log.Panic(err)
		replyMessage.Text = "There was an error saving your note"
	} else {
		replyMessage.Text = fmt.Sprintf("Note saved!\nUse `/get %s` to retrieve this note", tag)
		replyMessage.ParseMode = "markdown"
	}
	go bot.Send(replyMessage)
}

//GetNotes is used to get all notes of a subject
func GetNotes(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	reg := regexp.MustCompile("#\\w+")
	tag := reg.FindString(message.Text)

	log.Println(tag)
	if tag == "" {
		//Handle this
		return
	}

	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	notesCollection := client.Database("bot").Collection("notes")

	var notes []Note

	notesFromDb, err := notesCollection.Find(ctx, bson.M{"tag": bson.D{primitive.E{Key: "$eq", Value: tag}}})
	if err != nil {
		log.Panic("Error in finding notes", err)
		return
	}
	err = notesFromDb.All(ctx, &notes)
	if err != nil {
		log.Panic("Error in getting notes from DB", err)
		return
	}
	if len(notes) == 0 {
		noNotesMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "No notes were found for this tag")
		go bot.Send(noNotesMessage)
		return
	}
	for _, note := range notes {
		if note.FileID != "" {
			documentShare := tgbotapi.NewDocumentShare(update.Message.Chat.ID, note.FileID)
			documentShare.DisableNotification = true
			go bot.Send(documentShare)
		} else {
			forwardConfig := tgbotapi.NewForward(update.Message.Chat.ID, note.MessageFromChatID, note.MessageID)
			go bot.Send(forwardConfig)
		}
	}
}

//DeleteNote deletes a note
func DeleteNote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	//TODO() still in progress
	repliedToDocument := update.Message.ReplyToMessage.Document
	repliedToMessage := update.Message.ReplyToMessage
	if repliedToDocument == nil && repliedToMessage == nil {
		//Handle null document
		return
	}

	client, err := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	notesCollection := client.Database("bot").Collection("notes")
	var res *mongo.DeleteResult
	if repliedToDocument != nil {
		log.Println(repliedToDocument.FileID)
		res, err = notesCollection.DeleteOne(ctx, bson.M{"content": repliedToDocument.FileName})
	} else {
		log.Println(repliedToMessage.MessageID)
		res, err = notesCollection.DeleteOne(ctx, bson.M{"content": repliedToMessage.Text})
	}
	if err != nil {
		log.Panic(err)
	} else {
		log.Println(res)
		replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if res.DeletedCount == 0 {
			replyMessage.Text = "No notes were deleted"
		} else {
			replyMessage.Text = "Note was deleted"
		}
		go bot.Send(replyMessage)
	}
}
