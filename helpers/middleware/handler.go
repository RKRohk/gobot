package middleware

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	Start(bot *tgbotapi.BotAPI, update *tgbotapi.Update)
	Continue(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *UserSession)
	Done(bot *tgbotapi.BotAPI, update *tgbotapi.Update, session *UserSession)
}
