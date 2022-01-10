package reminders

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//ReminderInterrupt is an abstract type which is used to denote an interrupt
//An interrupt of this type is sent when a new reminder is closer than the current closest reminder
type reminderInterrupt struct {
}

//interruptChannel is a channel of reminderInterrupts
var interruptChannel = make(chan reminderInterrupt)

//RemindersChannel
var RemindersChannel = make(chan *EventToBeSent)

var reminders Reminders

//ClosestEvent represents a reminder which is closest to the current time
var ClosestEvent Reminder

func init() {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	time.Local = loc
	//Gets the closest reminder
	GetClosestReminder()
	go ReminderService()
}

//GetClosestReminder fetches the closest reminder from the database
func GetClosestReminder() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cursor, err := database.Client.Database("bot").Collection("reminders").Aggregate(ctx, mongo.Pipeline{
		{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "date", Value: 1}}}},
		{primitive.E{Key: "$limit", Value: 1}}})

	if err != nil {
		log.Println("Error fetching reminder from database", err)
		return
	}

	err = cursor.All(ctx, &reminders)
	if err != nil {
		log.Println("Error", err)
	}
	log.Println(reminders)
	if len(reminders) == 0 {
		log.Println("No reminders")
		ClosestEvent = Reminder{Date: time.Date(2050, 12, 31, 12, 12, 12, 12, time.UTC)}
	} else {
		ClosestEvent = reminders[0]
	}

}

//ReminderService runs for eternity till the bot stops
func ReminderService() {
	for {
		log.Println(ClosestEvent)
		select {
		case <-interruptChannel:
			{
				log.Println("I have got a new reminder")
				GetClosestReminder()
			}
		case <-time.After(time.Until(ClosestEvent.Date)):
			{
				log.Println("Say what?")
				SendReminder(&ClosestEvent)
				GetClosestReminder()
			}
		}

	}
}

func SendReminder(event *Reminder) {

	newMessage := tgbotapi.NewMessage(event.ChatId, "")

	newMessage.Text = fmt.Sprintf("Here is your reminder\n\n*%s*", strings.Trim(event.Title, ""))

	newMessage.ParseMode = "markdown"

	eventToBeSent := &EventToBeSent{Reminder: event, MessageConfig: &newMessage}

	RemindersChannel <- eventToBeSent

	event.Delete()
}
