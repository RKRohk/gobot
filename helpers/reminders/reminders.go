package reminders

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = log.New(os.Stdout, "Package:Reminders", log.LstdFlags)

//Reminder is memory representation of a reminder object
type Reminder struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Date   time.Time          `bson:"date,omitempty"`
	Title  string             `bson:"title,omitempty"`
	ChatId int64              `bson:"chatID,omitempty"`
}

type EventToBeSent struct {
	Reminder      *Reminder
	MessageConfig *tgbotapi.MessageConfig
}

//Reminders is an array of reminders
type Reminders []Reminder

//ParseDate parses given date (string) to a time.Time object
func ParseDate(date string) (time.Time, error) {
	format := "_2/1/2006 3:04PM MST"
	log.Println("Parsing date", date)
	return time.Parse(format, date)
}

//Save saves a reminder to the database
func (reminder *Reminder) Save() error {
	remindersCollection := database.Client.Database("bot").Collection("reminders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := remindersCollection.InsertOne(ctx, &reminder)
	logger.Println(res.InsertedID)
	if err == nil {
		CompareReminders(reminder)
	}
	return err
}

func (reminder *Reminder) Delete() error {
	remindersCollection := database.Client.Database("bot").Collection("reminders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := remindersCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: reminder.ID}})

	return err
}

func CompareReminders(reminder *Reminder) {
	if reminder.Date.Before(ClosestEvent.Date) {
		interruptChannel <- reminderInterrupt{}
	}
}
