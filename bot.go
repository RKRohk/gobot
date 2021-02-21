package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/rkrohk/gobot/handler"
	"github.com/rkrohk/gobot/helpers"

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

	go func() {
		updates, err := bot.GetUpdatesChan(u)

		if err != nil {
			log.Fatal("There was an error getting updates")
		}
		for update := range updates {
			if update.Message != nil { // ignore any non-Message Updates
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				if update.Message.Command() != "" {
					handler.Commandhandler(bot, update)
				}
			}

			if update.Message.Command() == "" && update.Message != nil {
				if update.Message.Chat.IsPrivate() || strings.Contains(update.Message.Text, "Rohk") || (update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.FirstName == bot.Self.FirstName) {
					helpers.GetMessage(bot, &update)
				}
			}

			if query := update.CallbackQuery; query != nil {
				log.Printf("CallbackQuery %s", query.Data)
				handler.CallbackQueryHandler(bot, &update)
			}

			if update.InlineQuery != nil {
				handler.Inlinehandler(bot, &update)
			}

		}
	}()

	//Added graceful shutdown
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, os.Kill)

	s := <-exitChan
	log.Println("Gracefully Shutting down the bot", s)

	bot.StopReceivingUpdates()
	return
}
