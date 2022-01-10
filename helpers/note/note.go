package note

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/database"
	"github.com/rkrohk/gobot/helpers/search"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

var notesCollection = database.Client.Database("bot").Collection("notes")

//BulkSaveNotes inserts the entire array of notes from
//the user session to the database
func BulkSaveNotes(notes []*Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	notesToBeInserted := make([]interface{}, len(notes))
	for i, v := range notes {
		notesToBeInserted[i] = *v
	}
	_, err := notesCollection.InsertMany(ctx, notesToBeInserted)
	log.Println(notesToBeInserted)
	return err
}

func BulkIndexNotes(bot *tgbotapi.BotAPI, notes []*Note) {
	for _, v := range notes {
		fileConfig := tgbotapi.FileConfig{FileID: v.FileID}
		doc, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Println("Note.Go: Error getting file link", err)
		}
		link := doc.Link(os.Getenv("BOT_TOKEN"))
		go search.IndexBulk(link, v.Content, v.FileID, v.Tag)
	}
}

func DeleteTag(tag string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	deleteConfig := bson.D{primitive.E{Key: "tag", Value: tag}}
	res, err := notesCollection.DeleteMany(ctx, deleteConfig)
	if err != nil || res.DeletedCount == 0 {
		return 0, err
	}
	return res.DeletedCount, err
}
