package helpers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//ParseDate parses string to date
func ParseDate(inputString string) (time.Time, error) {
	format := "02/01/2006 3:04 PM MST"
	date, err := time.Parse(format, inputString)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return date, err
	}
	return date, nil

}

//Timer starts a timer
func Timer(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := update.Message

	//Regex to extract time
	reg := regexp.MustCompile("([0-9]{1,2})/([0-9]{1,2})/202([0-9]) ([1-2]?)([0-9]):([0-9]{2}) (A|P)M ([A-Z])ST")

	remindTimeStr := reg.FindString(msg.Text)
	remindTimeParsed, err := ParseDate(remindTimeStr)
	if err != nil {
		log.Println("Invalid time")
		log.Println(remindTimeStr)
		return
	}

	//Removing time and command from the message
	remindMessage := reg.ReplaceAllString(msg.Text, "")
	remindMessage = strings.Replace(remindMessage, "/timer", "", 1)

	//Making and parsing the reply message
	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Okay, I will remind you about:\n*%s*\n\nAfter `%s`", remindMessage, remindTimeStr))
	message.ParseMode = "markdown"
	message.ReplyToMessageID = update.Message.MessageID
	waitTime := remindTimeParsed.Sub(time.Now())

	//Checking if the date given is before current time
	if waitTime < 0 {
		message.Text = "I cannot remind you in the past"
		bot.Send(message)
		return
	}
	sentMessage, err := bot.Send(message)

	//Bot sleeps till remind time
	time.Sleep(waitTime)
	message.Text = fmt.Sprintf("Yo here is your reminder\n*%s*", remindMessage)
	bot.Send(message)

	//Deleting the previous message
	deleteOldMessage := tgbotapi.NewDeleteMessage(sentMessage.Chat.ID, sentMessage.MessageID)
	bot.Send(deleteOldMessage)
}
