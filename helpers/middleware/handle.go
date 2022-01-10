package middleware

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Handle(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	currentSession := session[UserChatID(update.Message.Chat.ID)]

	switch *currentSession.CurrentCommand {
	case "batchsave":
		{
			//TODO()
			middlwareLogger.Println("Command is batchsave")
			currentSession.SessionHandler.Continue(bot, update, currentSession)
			return
		}
	case "done":
		{
			currentSession.SessionHandler.Done(bot, update, currentSession)
			ClearSession(update)
			break
		}
	case "cancel":
		{
			ClearSession(update)
			break
		}

	}

}
