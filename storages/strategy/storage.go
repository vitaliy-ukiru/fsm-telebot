package strategy

import "github.com/vitaliy-ukiru/fsm-telebot"

// Strategy for addressing. It works as bit set.
//
// So far, you can only combine [User] and [Chat],
// but with the support of forums in telebot, we plan to add them as well.
//
// Have "magic" value - [Empty]. It needs for safe support zero value.
// This value must be just pass, without any logic.
type Strategy byte

// Empty strategy is unset value, but it supports like [Default] strategy
const Empty Strategy = 0

const (
	_ Strategy = 1 << iota

	// User addressing. It will make one state for one user in all chats.
	User

	// Chat addressing. It will make one state for all users in one chat.
	Chat

	// Default is contains user and chat addressing.
	// It will make state for every user in every chat.
	Default = User | Chat
)

func (s Strategy) String() string {
	switch s {
	case Empty:
		return "strategy.Empty"
	case Default:
		return "strategy.Default"
	case User:
		return "strategy.OnlyUser"
	case Chat:
		return "strategy.OnlyChat"
	}
	return "strategy.INVALID"
}

// Storage works over base storage and applies strategy.
type Storage struct {
	storage  fsm.Storage
	strategy Strategy
}

func NewStorage(storage fsm.Storage, strategy Strategy) *Storage {
	return &Storage{storage: storage, strategy: strategy}
}

func (s *Storage) Strategy() Strategy {
	return s.strategy
}

func (s *Storage) SetStrategy(strategy Strategy) {
	s.strategy = strategy
}

func (s *Storage) GetState(c, u int64) (fsm.State, error) {
	c, u = s.strategy.apply(c, u)
	return s.storage.GetState(c, u)
}

func (s *Storage) SetState(c, u int64, state fsm.State) error {
	c, u = s.strategy.apply(c, u)
	return s.storage.SetState(c, u, state)
}

func (s *Storage) ResetState(c, u int64, withData bool) error {
	c, u = s.strategy.apply(c, u)
	return s.storage.ResetState(c, u, withData)
}

func (s *Storage) UpdateData(c, u int64, key string, data any) error {
	c, u = s.strategy.apply(c, u)
	return s.storage.UpdateData(c, u, key, data)
}

func (s *Storage) GetData(c, u int64, key string, to any) error {
	c, u = s.strategy.apply(c, u)
	return s.storage.GetData(c, u, key, to)
}

func (s *Storage) Close() error {
	return s.storage.Close()
}

func (s Strategy) apply(chat, user int64) (int64, int64) {
	if s == Empty {
		return chat, user
	}

	if s&Chat == 0 {
		chat = 0
	}

	if s&User == 0 {
		user = 0
	}

	return chat, user
}
