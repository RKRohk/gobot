package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//APIMessage represents a response from API
type APIMessage struct {
	Message string `json:"message,omitempty"`
}

//SendMessage sends a message to the bot to store
func SendMessage(message string) {
	base, err := url.Parse("http://chatbotapi/setmessage")
	if err != nil {
		log.Fatal("Invalid API URI")
	}
	params := url.Values{}
	params.Add("message", message)
	base.RawQuery = params.Encode()

	// fmt.Printf("Encoded URL is %q\n", base.String())
	_, err = http.Get(base.String())
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer resp.Body.Close()
	// var target Message
	// err = json.NewDecoder(resp.Body).Decode(&target)
	// if err != nil {
	// 	log.Panic(err)
	// } else {
	// 	log.Println(target.Message)
	// }
}

//GetMessage gets message from the bot and replies to the user
func GetMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	base, err := url.Parse("http://chatbotapi/getmessage")
	params := url.Values{}
	params.Add("message", message.Text)
	base.RawQuery = params.Encode()

	// fmt.Printf("Encoded URL is %q\n", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var target APIMessage
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		log.Panic(err)
	} else {
		replyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, target.Message)
		replyMessage.ReplyToMessageID = message.MessageID
		bot.Send(replyMessage)
	}
}
