package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/rkrohk/gobot/handler"
	"github.com/rkrohk/gobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Unable to read token", err)
	}
	var blockedUser int
	if blockedUser, err = strconv.Atoi(os.Getenv("BLOCKED_USER")); err != nil {
		log.Panic("Unable to read blocked user")
	} else {
		log.Println("Blocked user is", blockedUser)
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
				} else if update.Message.Command() == "" && update.Message != nil && update.Message.Text != "" {
					if update.Message.Chat.IsPrivate() || strings.Contains(update.Message.Text, "Gora") || (update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.FirstName == bot.Self.FirstName) {
						if update.Message.From.ID == blockedUser {

						} else {
							helpers.GetMessage(bot, &update)
						}

					} else {
						helpers.SendMessage(update.Message.Text)
					}
				}
			} else if query := update.CallbackQuery; query != nil {
				log.Printf("CallbackQuery %s", query.Data)
				handler.CallbackQueryHandler(bot, &update)
			} else if update.InlineQuery != nil {
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
