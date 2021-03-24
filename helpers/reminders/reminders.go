package reminders

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/database"
)

var logger = log.New(os.Stdout, "Package:Reminders", log.LstdFlags)

type Reminder struct {
	ID        int              `bson:"_id,omitempty"`
	Message   tgbotapi.Message `bson:"message,omitempty"`
	Date      time.Time        `bson:"date,omitempty"`
	Title     string           `bson:"title,omitempty"`
	CreatedBy int              `bson:"created_by,omitempty"`
}

type Reminders []Reminder

func ParseDate(date string) (time.Time, error) {
	format := "_2/1/2006 3:04PM"
	return time.Parse(format, date)
}

func (reminder *Reminder) Save() error {
	remindersCollection := database.Client.Database("bot").Collection("reminders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := remindersCollection.InsertOne(ctx, &reminder)
	logger.Println(res.InsertedID)
	return err
}
