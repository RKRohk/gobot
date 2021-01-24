package helpers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Slap takes two users as argument and slaps one
func Slap(user1 string, user2 string) string {
	pick, err := GetSlapStrings()
	if err != nil {
		return "database error occurred"
	}
	return fmt.Sprintf(pick, user1, user2)
}

//ShowSlapStrings returns widgets showing all slap strings and giving user an option to delete any
func ShowSlapStrings(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here are all the slap strings")
	slaps, err := GetAllSlapStrings()
	if err != nil {
		log.Panic(err)
		msg.Text = "There was an error"
		bot.Send(msg)
		return
	}
	replies := make([][]tgbotapi.InlineKeyboardButton, len(slaps))
	for i, v := range slaps {
		callbackData := fmt.Sprintf("deleteslap:%s", v.ID.Hex())
		replies[i] = append(replies[i], tgbotapi.InlineKeyboardButton{Text: v.Text, CallbackData: &callbackData})
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: replies}

	//TODO() add a cancel button for this
	msg.ReplyMarkup = replyMarkup
	bot.Send(msg)

}

//DeleteSlapString deletes a slap string from database
func DeleteSlapString(bot *tgbotapi.BotAPI, update *tgbotapi.Update, documentID string) {
	editedMessage := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
	if err := DeleteSlapFromDb(documentID); err != nil {
		editedMessage.Text = "There was some problem deleting the slap string. Please try later"
		editedMessage.ReplyMarkup = nil
		go bot.Send(editedMessage)
		return
	}
	editedMessage.Text = "Slap String deleted from database"
	go bot.Send(editedMessage)
	editedMessage.ReplyMarkup = nil

}
