package fsm

import (
	tele "gopkg.in/telebot.v3"
)

type Strategy int

const (
	StrategyUserInChat Strategy = iota
	StrategyChat
	StrategyGlobalUser
	StrategyUserInTopic
	StrategyChatTopic

	StrategyDefault = StrategyUserInChat
)

func (s Strategy) Apply(botId int64, chatId int64, userId int64, threadId int64) StorageKey {
	switch s {
	case StrategyChat:
		return StorageKey{
			BotID:  botId,
			ChatID: chatId,
			UserID: chatId,
		}
	case StrategyGlobalUser:
		return StorageKey{
			BotID:  botId,
			ChatID: userId,
			UserID: userId,
		}
	case StrategyUserInTopic:
		return StorageKey{
			BotID:    botId,
			ChatID:   chatId,
			UserID:   userId,
			ThreadID: threadId,
		}
	case StrategyChatTopic:
		return StorageKey{
			BotID:    botId,
			ChatID:   chatId,
			UserID:   chatId,
			ThreadID: threadId,
		}
	case StrategyUserInChat:
		fallthrough
	default:
		return StorageKey{
			BotID:  botId,
			ChatID: chatId,
			UserID: userId,
		}
	}
}

func ExtractKeyWithStrategy(c tele.Context, strategy Strategy) StorageKey {
	chatId := c.Chat().ID
	userId := c.Sender().ID
	threadId := int64(c.Message().ThreadID) // TODO: check is safe using zero thread id

	bot := c.Bot().Me
	return strategy.Apply(bot.ID, chatId, userId, threadId)
}
