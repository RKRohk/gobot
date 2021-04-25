package middleware

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//UserChatID refers to telegram UserChatID
type UserChatID int64

//Session denotes a session that a user has with the bot
type UserSession struct {
	ChatID         UserChatID
	CurrentCommand *string
	SessionHandler Handler
}

//Session is a map of int64 (i.e ChatID) and UserSession struct
type Session map[UserChatID]*UserSession

var session = make(Session)

//Done function removes the current command from the user's session
func Done(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	session := GetSession(update)
	session.SessionHandler.Done(bot, update, session)
	ClearSession(update)
	return nil
}

func Cancel(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	ClearSession(update)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Operation Cancelled...")
	bot.Send(reply)
}

//HasSession checks whether the user has a current command pending in the session
func HasSession(update *tgbotapi.Update) bool {
	return session[UserChatID(update.Message.Chat.ID)] != nil
}

func GetSession(update *tgbotapi.Update) *UserSession {
	return session[UserChatID(update.Message.Chat.ID)]
}

func AddSession(chatID UserChatID, currentCommand string, handler Handler) {
	session[chatID] = &UserSession{ChatID: chatID, CurrentCommand: &currentCommand, SessionHandler: handler}
}

func ClearSession(update *tgbotapi.Update) {
	delete(session, UserChatID(update.Message.Chat.ID))
	//TODO()
	middlwareLogger.Println("Session Deleted")
}
