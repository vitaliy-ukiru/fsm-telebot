package fsm_telebot

import (
	tele "gopkg.in/telebot.v3"
)

// FSMContext is wrapper for work with FSM from handlers and related to telebot.Context.
type FSMContext interface {
	// Bot returns the bot instance.
	Bot() *tele.Bot

	// State returns current state for sender.
	State() State

	// Set state for sender.
	Set(state State)

	// Finish state for sender and deletes data if set true.
	Finish(deleteData bool) error

	// Update data in storage.
	Update(key string, data interface{}) error
	// Get data from storage.
	// Data will be nil if storage not contains object with given key and error will be ErrNotFound
	Get(key string) (data interface{}, err error)

	// MustGet returns data from storage and ignore any errors.
	MustGet(key string) (data interface{})
}

type fsmContext struct {
	s          Storage
	c          tele.Context
	chat, user int64
}

// NewFSMContext returns new FSMContext
func NewFSMContext(c tele.Context, storage Storage) FSMContext {
	return &fsmContext{
		c:    c,
		s:    storage,
		chat: c.Chat().ID,
		user: c.Sender().ID,
	}
}

func (f *fsmContext) Bot() *tele.Bot {
	return f.c.Bot()
}

func (f *fsmContext) State() State {
	return f.s.GetState(f.chat, f.user)
}

func (f *fsmContext) Set(state State) {
	f.s.SetState(f.chat, f.user, state)
}

func (f *fsmContext) Finish(deleteData bool) error {
	return f.s.ResetState(f.chat, f.user, deleteData)
}

func (f *fsmContext) Update(key string, data interface{}) error {
	return f.s.UpdateData(f.chat, f.user, key, data)
}

func (f *fsmContext) Get(key string) (interface{}, error) {
	return f.s.GetData(f.chat, f.user, key)
}

func (f *fsmContext) MustGet(key string) interface{} {
	data, _ := f.s.GetData(f.chat, f.user, key)
	return data
}
