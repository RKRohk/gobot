package helpers

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type reminder struct {
	text         string
	chatID       int64
	reminderTime time.Time
}

func newReminder(text string, chatID int64, reminderTime string) *reminder {

	return &reminder{text, chatID, time.Now()}
}

//Remind takes an update as argument and handles the function call
func Remind(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := update.Message

	msgArr := strings.Split(msg.Text, " ")

	if len(msgArr) == 0 {
		replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid message")
		replyMessage.ReplyToMessageID = update.Message.MessageID
		bot.Send(replyMessage)
		return
	}
	remindTimeStr := msgArr[len(msgArr)-1]

	if waitDuration, err := time.ParseDuration(remindTimeStr); err != nil {
		replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid Time")
		replyMessage.ReplyToMessageID = update.Message.MessageID
		bot.Send(replyMessage)
	} else {
		remindMessage := strings.Join(msgArr[1:len(msgArr)-1], " ")
		fmt.Println(msgArr[1 : len(msgArr)-1])
		fmt.Println()
		message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Okay, I will remind you about:\n*%s*\n\nAfter `%s`", remindMessage, remindTimeStr))
		message.ParseMode = "markdown"
		message.ReplyToMessageID = update.Message.MessageID
		bot.Send(message)
		time.Sleep(waitDuration)
		message.Text = fmt.Sprintf("Yo here is your reminder\n*%s*", remindMessage)
		bot.Send(message)
	}

	fmt.Println(remindTimeStr)
}
