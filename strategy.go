package fsm

import (
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/v2/internal/null"
	tele "gopkg.in/telebot.v3"
)

type Strategy int

const (
	StrategyUserInChat Strategy = iota
	StrategyChat
	StrategyGlobalUser
	StrategyUserInTopic
	StrategyChatTopic
	_maxStrategy

	StrategyDefault = StrategyUserInChat
)

// explicit set type for check in compile helps
// check if items in array don't equal actually count of values
//
//goland:noinspection GoVarAndConstTypeMayBeOmitted
var strategyStr [_maxStrategy]string = [...]string{
	StrategyUserInChat:  "StrategyUserInChat",
	StrategyChat:        "StrategyChat",
	StrategyGlobalUser:  "StrategyGlobalUser",
	StrategyUserInTopic: "StrategyUserInTopic",
	StrategyChatTopic:   "StrategyChatTopic",
}

func (s Strategy) String() string {
	if 0 <= s && s <= _maxStrategy {
		return strategyStr[s]
	}
	return fmt.Sprintf("Strategy(%d)", s)
}

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

func extractKeyWithStrategy(c tele.Context, strategy Strategy) (StorageKey, bool) {
	var (
		chatId null.Nullable[int64]
		userId null.Nullable[int64]
	)

	if chat := c.Chat(); chat != nil {
		chatId.Set(chat.ID)
	}

	if user := c.Sender(); user != nil {
		userId.Set(user.ID)
	}

	var threadId int64
	if m := c.Message(); m != nil && m.TopicMessage {
		threadId = int64(c.Message().ThreadID)
	}

	if !chatId.Valid {
		chatId = userId
	}

	if !chatId.Valid || !userId.Valid {
		return StorageKey{}, false
	}

	bot := c.Bot().Me
	return strategy.Apply(bot.ID, chatId.Value, userId.Value, threadId), true
}
