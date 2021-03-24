package helpers

import (
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers/reminders"
)

//Remind handles the /remind command
func Remind(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	text := message.Text
	text = strings.Replace(text, "/remind ", "", 1)

	log.Println(text)
	dateRegex := regexp.MustCompile(`0?\d\/0?\d\/20\d\d 0?\d:\d\d[A|P]M`)

	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "hi")
	reply.ReplyToMessageID = message.From.ID

	defer func() {
		go bot.Send(reply)
	}()

	if dateIndices := dateRegex.FindStringIndex(text); dateIndices != nil {
		dateString := text[dateIndices[0]:dateIndices[1]]
		log.Println(dateString, ":DateString")
		if date, err := reminders.ParseDate(dateString); err != nil {
			reply.Text = "Time format is incorrect"
			log.Println("Invalid date format", err)
			return
		} else {
			log.Println("Parsed date is ", date)
			title := text[dateIndices[1]:]

			newReminder := &reminders.Reminder{Date: date, Title: title, ChatId: message.Chat.ID}
			err := newReminder.Save()
			if err != nil {
				reply.Text = "There was an error saving your reminder"
				log.Println(err)
			} else {
				reply.Text = "Reminder saved"
				log.Println("Reminder saved")
			}
		}

	} else {
		log.Println(dateIndices)
		log.Println("Invalid date")
		reply.Text = "Please enter a valid date"
	}

}
