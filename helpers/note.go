package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/search"
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

//Extracts the hash tag from the message
func ExtractTag(message string) string {
	reg := regexp.MustCompile("#\\w+")
	return reg.FindString(message)
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

	tag := ExtractTag(message.Text)

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
		fileConfig := tgbotapi.FileConfig{FileID: repliedToDocument.FileID}
		doc, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Println("Note.Go: Error getting file link", err)
		}
		link := doc.Link(os.Getenv("BOT_TOKEN"))
		go search.Index(link, repliedToDocument, tag)
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
		if _, err := bot.Send(noNotesMessage); err != nil {
			log.Println("Error message ", tag, "\n", err)
		}
		return
	}
	for _, note := range notes {
		if note.FileID != "" {
			documentShare := tgbotapi.NewDocumentShare(update.Message.Chat.ID, note.FileID)
			documentShare.DisableNotification = true
			if _, err := bot.Send(documentShare); err != nil {
				log.Println("Error sending file for ", tag, "\n", err)
			}
		} else {
			forwardConfig := tgbotapi.NewForward(update.Message.Chat.ID, note.MessageFromChatID, note.MessageID)
			if _, err := bot.Send(forwardConfig); err != nil {
				log.Println("Error forwarding message for ", tag, "\n", err)
			}
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

func GetAllTags(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

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

	distinctTags, err := notesCollection.Distinct(context.Background(), "tag", bson.D{})

	returnMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	returnMessage.ParseMode = "markdown"

	if err != nil {
		log.Println("Error", err)
		returnMessage.Text = "There was an error finding your tags"
	} else {
		log.Println(distinctTags)
		returnMessageText := "*Your saved notes are*\n\n"
		for _, tag := range distinctTags {
			returnMessageText = returnMessageText + fmt.Sprintf("• %s\n", tag)
		}
		fmt.Println(returnMessageText)
		returnMessage.Text = returnMessageText
	}
	bot.Send(returnMessage)

}
