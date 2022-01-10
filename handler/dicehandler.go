package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type DiceOption struct {
	string
	int
}

func DiceHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	dice := update.Message.Dice

	emoji := dice.Emoji
	value := dice.Value
	diceOption := DiceOption{emoji, value}

	switch diceOption {
	case DiceOption{"🎲", 6}:
		fallthrough
	case DiceOption{"🎯", 6}:
		fallthrough
	case DiceOption{"🎳", 6}:
		fallthrough
	case DiceOption{"🏀", 5}:
		fallthrough
	case DiceOption{"⚽", 5}:
		fallthrough
	case DiceOption{"🎰", 64}:
		{
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Pro!")
			reply.ReplyToMessageID = update.Message.MessageID
			bot.Send(reply)
		}
	}
}
