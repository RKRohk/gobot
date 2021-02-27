package handler

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers"
)

//Returns slap string for inline "slap" query
func inlineSlap(update *tgbotapi.Update) string {
	log.Println(update.InlineQuery.Query)
	currentUser := update.InlineQuery.From.FirstName
	victimUser := ""
	messageWithoutCommand := strings.Replace(update.InlineQuery.Query, "slap ", "", 1)
	victimUser = messageWithoutCommand
	slaptext := helpers.Slap(currentUser, victimUser)

	return slaptext
}

//Inlinehandler handles inline queries
func Inlinehandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

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

	slapArticle := tgbotapi.NewInlineQueryResultArticleMarkdown("2", "Slap", inlineSlap(update))
	article.Description = "¯\\_(ツ)_/¯"
	slapArticle.Description = "Write the name of the person you want to slap"

	inlineConf = tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       []interface{}{article, slapArticle},
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println("Error sending slap inline query answer")
		log.Println(err)
	}
}
