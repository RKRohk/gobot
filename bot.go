package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/rkrohk/gobot/database"
	"github.com/rkrohk/gobot/handler"
	"github.com/rkrohk/gobot/helpers"
	"github.com/rkrohk/gobot/helpers/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI
var err error
var blockedUser int
var owner int

func init() {
	token := os.Getenv("BOT_TOKEN")
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Unable to read token", err)
	}

	if owner, err = strconv.Atoi(os.Getenv("OWNER")); err != nil {
		log.Fatal("Unable to read owner user")
	}

	if blockedUser, err = strconv.Atoi(os.Getenv("BLOCKED_USER")); err != nil {
		log.Fatal("Unable to read blocked user")
	} else {
		log.Println("Blocked user is", blockedUser)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//Initializing Reminder service
	log.Println("Initializing Reminder service")
	go helpers.InitReminderService(bot)

}

func teardown() {
	database.DisconnectDatabase()
	bot.StopReceivingUpdates()
	ai.Close()
}

func main() {

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

	defer func() {
		if err := recover(); err != nil {
			log.Println("Error occurred", err)
			errorMessage := tgbotapi.NewMessage(int64(owner), "An error occurred")
			bot.Send(errorMessage)
			log.Println("Recovering")
		}
	}()

	//Added graceful shutdown
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, syscall.SIGTERM)

	s := <-exitChan
	log.Println("Gracefully Shutting down the bot", s)

	teardown()

}
