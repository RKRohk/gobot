package reminder

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Reminder struct {
	ID      int              `bson:"_id,omitempty"`
	Message tgbotapi.Message `bson:"message,omitempty"`
	Date    time.Time        `bson:"date,omitempty"`
	title   string           `bson:"title,omitempty"`
}

func CreateReminder(date time.Time, title string, message tgbotapi.Message) *Reminder {
	return nil
}
