package helpers

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//ParseDate parses string to date
func ParseDate(inputString string) (time.Time, error) {
	format := "2006-01-02 15:04 MST"
	fmt.Println(inputString)

	date, err := time.Parse(format, inputString)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return date, err
	}
	fmt.Println("Date is")
	fmt.Println(date)
	return date, nil

}

//Timer starts a timer
func Timer(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	//Removes command from message
	timeStr := strings.Replace(update.Message.Text, "/timer ", "", 1)
	//Removing whitespaces from message
	timeStr = strings.Trim(timeStr, " ")
	//Converting text to int with base 10
	endTime, err := ParseDate(timeStr)
	startTime := time.Now()
	waitTimeFloat := endTime.Sub(startTime).Seconds()
	waitTime := int(waitTimeFloat)

	//If there is any error parsing the int, return
	if err != nil {
		fmt.Println(err)
		return
	}

	//Sending first message
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", "Starting Countdown"))
	//Creating a new inline keyboard button for the "Cancel" operation

	uuid := GenerateUUID()
	button := tgbotapi.NewInlineKeyboardButtonData("Cancel", uuid)
	buttons := make([]tgbotapi.InlineKeyboardButton, 1, 1)
	buttons[0] = button
	menu := tgbotapi.NewInlineKeyboardMarkup(buttons)

	//First message (This message has to be referenced for all editing operations)
	sentMessage, err := bot.Send(msg)
	sentMessage.Command()
	for i := waitTime; i >= 0; i-- {
		//First message is edited at every iteration of the loop
		var editedMessage tgbotapi.EditMessageTextConfig

		if i == 0 {
			editedMessage = tgbotapi.NewEditMessageText(sentMessage.Chat.ID, sentMessage.MessageID, "Timer complete")
		} else {
			editedMessage = tgbotapi.NewEditMessageText(sentMessage.Chat.ID, sentMessage.MessageID, fmt.Sprintf("%d", i))
			editedMessage.ReplyMarkup = &menu
		}
		go bot.Send(editedMessage)
		//Sleep for one second
		time.Sleep(time.Second)
	}
}
