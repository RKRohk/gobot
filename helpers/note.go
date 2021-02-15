package helpers

import (
	"context"
	"log"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Note denotes how a note is saved in the database
type Note struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	FileID string             `bson:"fileID,omitempty"`
	Tag    string             `bson:"tag,omitempty"`
}

//SaveNote saves a note
func SaveNote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	repliedToDocument := update.Message.ReplyToMessage.Document
	if repliedToDocument == nil {
		//Handle null document
		return
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
	newNote := Note{FileID: repliedToDocument.FileID, Tag: tag}
	log.Println(newNote)
	_, err = notesCollection.InsertOne(ctx, &newNote)
	if err != nil {
		log.Panic(err)
	}
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

	notesFromDb, err := notesCollection.Find(ctx, bson.M{"tag": bson.D{{"$eq", tag}}})
	if err != nil {
		log.Panic("Error in finding notes", err)
		return
	}
	err = notesFromDb.All(ctx, &notes)
	if err != nil {
		log.Panic("Error in getting notes from DB", err)
		return
	}
	log.Println(notes)
	for _, note := range notes {
		documentShare := tgbotapi.NewDocumentShare(update.Message.Chat.ID, note.FileID)
		go bot.Send(documentShare)
	}
}

//DeleteNote deletes a note
func DeleteNote(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	repliedToDocument := update.Message.ReplyToMessage.Document
	if repliedToDocument == nil {
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
	_, err = notesCollection.DeleteOne(ctx, bson.M{"fileID": bson.D{{"$eq", repliedToDocument.FileID}}})
	if err != nil {
		log.Panic(err)
	}
}
