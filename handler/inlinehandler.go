package handler

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rkrohk/gobot/helpers"
	"github.com/rkrohk/gobot/helpers/slap"
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

	var articles []interface{}

	query := update.InlineQuery.Query
	query = strings.Trim(query, " ")

	if query == "" {
		//Shrug Article
		article := tgbotapi.NewInlineQueryResultArticleHTML("Shrug", "¯\\_(ツ)_/¯", fmt.Sprintf("¯\\_(ツ)_/¯"))
		article.Description = "¯\\_(ツ)_/¯"

		articles = append(articles, article)

		slapArticle := tgbotapi.NewInlineQueryResultArticleHTML("Slap", "Slap", "Slap")
		slapArticle.Description = "Write the name of the person you want to slap"
		articles = append(articles, slapArticle)
	}

	if strings.HasPrefix(query, "slap") {
		slapArticles := slap.InlineSlap(update)
		for _, v := range slapArticles {
			articles = append(articles, v)
		}
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    false,
		CacheTime:     0,
		Results:       articles,
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println("Error sending inline query answer")
		log.Println(err)

	}
}
