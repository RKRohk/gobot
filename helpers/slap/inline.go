package slap

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rkrohk/gobot/helpers"
)

//Returns slap string for inline "slap" query
func makeSlaps(update *tgbotapi.Update) []string {
	log.Println(update.InlineQuery.Query)
	currentUser := update.InlineQuery.From.FirstName
	victimUser := ""
	messageWithoutCommand := strings.Replace(update.InlineQuery.Query, "slap ", "", 1)
	victimUser = messageWithoutCommand
	slapStringsFromDb, err := helpers.GetAllSlapStrings()
	if err != nil {
		log.Println(err)
		return []string{"There was an error fetching slaps"}

	}
	slapStrings := make([]string, len(slapStringsFromDb))
	for i, v := range slapStringsFromDb {
		slapStrings[i] = fmt.Sprintf(v.Text, currentUser, victimUser)
	}
	return slapStrings

}

func InlineSlap(update *tgbotapi.Update) []tgbotapi.InlineQueryResultArticle {
	slapStrings := makeSlaps(update)

	var articles []tgbotapi.InlineQueryResultArticle

	for i, v := range slapStrings {
		articles = append(articles, tgbotapi.NewInlineQueryResultArticleHTML(fmt.Sprintf("slap_%d", i), v, v))
	}

	return articles
}
