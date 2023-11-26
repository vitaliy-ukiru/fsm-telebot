package fsm

import tele "gopkg.in/telebot.v3"

type Strategy int

const (
	StrategyUserInChat Strategy = iota
	StrategyChat
	StrategyGlobalUser
	StrategyUserInTopic
	StrategyChatTopic
)

func (s Strategy) Apply(chatId int64, userId int64, threadId *int64) StorageKey {
	switch s {
	case StrategyChat:
		return StorageKey{ChatID: chatId, UserID: chatId}
	case StrategyGlobalUser:
		return StorageKey{ChatID: userId, UserID: userId}
	case StrategyUserInTopic:
		return StorageKey{ChatID: chatId, UserID: userId, ThreadID: threadId}
	case StrategyChatTopic:
		return StorageKey{ChatID: chatId, UserID: chatId, ThreadID: threadId}
	case StrategyUserInChat:
		fallthrough
	default:
		return StorageKey{ChatID: chatId, UserID: userId}
	}
}

func ExtractKeyWithStrategy(c tele.Context, strategy Strategy) StorageKey {
	chatId := c.Chat().ID
	userId := c.Sender().ID
	var threadId *int64
	if msg := c.Message(); msg.TopicMessage {
		thread := int64(msg.ThreadID)
		threadId = &thread
	}

	return strategy.Apply(chatId, userId, threadId)
}
