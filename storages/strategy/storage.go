package strategy

import "github.com/vitaliy-ukiru/fsm-telebot"

type Strategy int

const (
	// Default strategy is chat + user
	Default Strategy = iota
	OnlyChat
	OnlyUser
)

func (s Strategy) String() string {
	switch s {
	case Default:
		return "strategy.Default"
	case OnlyUser:
		return "strategy.OnlyUser"
	case OnlyChat:
		return "strategy.OnlyChat"
	}
	return "strategy.INVALID"
}

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
	switch s {
	case OnlyChat:
		return chat, 0
	case OnlyUser:
		return 0, user
	case Default:
		fallthrough
	default:
		return chat, user
	}
}
