package handler

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Inlinehandler handles inline queries
func Inlinehandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.InlineQuery.Query

	if message == "" || message == " " {
		return
	}

	article := tgbotapi.NewInlineQueryResultArticleHTML("1", "Shrug", fmt.Sprintf("¯\\_(ツ)_/¯"))

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println("Error sending inline query answer")
		log.Println(err)
	}
}
