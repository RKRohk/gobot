package handler

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/helpers"
)

//CallbackQueryHandler handles callbackqueries
func CallbackQueryHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	callbackQuery := update.CallbackQuery.Data
	log.Println(callbackQuery)

	if strings.HasPrefix(callbackQuery, "deleteslap:") {
		documentID := strings.Replace(callbackQuery, "deleteslap:", "", 1)
		log.Println("LOL")
		log.Println(documentID)
		go helpers.DeleteSlapString(bot, update, documentID)
	}
}
