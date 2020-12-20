package main

import (
	"log"
	"os"

	"./handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(token)
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Command() != "" {
				handler.Commandhandler(bot, update)
			}
		}

		if query := update.CallbackQuery; query != nil {
			log.Printf("CallbackQuery %s", query.Data)
		}

		if update.InlineQuery != nil {
			handler.Inlinehandler(bot, &update)
		}

	}
}
