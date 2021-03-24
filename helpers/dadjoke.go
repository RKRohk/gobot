package helpers

import (
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//DadJoke fetches a dadjoke from icanhazdadjoke.com
func DadJoke(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	url := "https://icanhazdadjoke.com"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "text/plain")

	res, err := http.DefaultClient.Do(req)

	var response = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if err != nil {
		response.Text = "No jokes for you!"
	} else {
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			response.Text = "No jokes for you!"
		} else {
			response.Text = string(body)
		}
	}

	bot.Send(response)
}
