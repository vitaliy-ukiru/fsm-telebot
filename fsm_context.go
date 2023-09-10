package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// Context is wrapper for work with FSM from handlers and related to telebot.Context.
type Context interface {
	// Bot returns the bot instance.
	Bot() *tele.Bot

	// State returns current state for sender.
	State() (State, error)

	// Set state for sender.
	Set(state State) error

	// Finish state for sender and deletes data if set true.
	Finish(deleteData bool) error

	// Update data in storage.
	Update(key string, data any) error

	// Get data from storage and save it to `to`.
	// `to` must be a pointer.
	// Data will be nil if storage not contains object with given key and error will be ErrNotFound
	Get(key string, to any) error

	// MustGet returns data from storage and save it to `to` ignoring errors.
	// `to` must be a pointer.
	MustGet(key string, to any)
}

type BuiltinContext struct {
	s          Storage
	c          tele.Context
	chat, user int64
}

// NewFSMContext returns new builtin FSM Context.
func NewFSMContext(c tele.Context, storage Storage) *BuiltinContext {
	return &BuiltinContext{
		c:    c,
		s:    storage,
		chat: c.Chat().ID,
		user: c.Sender().ID,
	}
}

func (f *BuiltinContext) Bot() *tele.Bot {
	return f.c.Bot()
}

func (f *BuiltinContext) State() (State, error) {
	return f.s.GetState(f.chat, f.user)
}

func (f *BuiltinContext) Set(state State) error {
	return f.s.SetState(f.chat, f.user, state)
}

func (f *BuiltinContext) Finish(deleteData bool) error {
	return f.s.ResetState(f.chat, f.user, deleteData)
}

func (f *BuiltinContext) Update(key string, data any) error {
	return f.s.UpdateData(f.chat, f.user, key, data)
}

func (f *BuiltinContext) Get(key string, to any) error {
	return f.s.GetData(f.chat, f.user, key, to)
}

func (f *BuiltinContext) MustGet(key string, to any) {
	_ = f.s.GetData(f.chat, f.user, key, to)
}
