package middleware

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//UserChatID refers to telegram UserChatID
type UserChatID int64

//Session denotes a session that a user has with the bot
type UserSession struct {
	ChatID         UserChatID
	CurrentCommand *string
}

//Session is a map of int64 (i.e ChatID) and UserSession struct
type Session map[UserChatID]UserSession

var session = make(Session)

//Done function removes the current command from the user's session
func (s *UserSession) Done() error {
	s.CurrentCommand = nil
	return nil
}

//HasSession checks whether the user has a current command pending in the session
func HasSession(update *tgbotapi.Update) bool {
	return session[UserChatID(update.Message.Chat.ID)].CurrentCommand != nil
}
