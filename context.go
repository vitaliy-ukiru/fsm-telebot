package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// Context is wrapper for work with FSM from handlers
// and related to telebot.Context.
type Context interface {
	// State returns current state for sender.
	State(ctx context.Context) (State, error)

	// SetState state for sender.
	SetState(ctx context.Context, state State) error

	// Finish state for sender and deletes data if arg provided.
	Finish(ctx context.Context, deleteData bool) error

	// Update data in storage. When data argument is nil it must
	// delete this item.
	Update(ctx context.Context, key string, data any) error

	// Data gets from storage and save it into `to` argument.
	// Destination argument must be a valid pointer.
	Data(ctx context.Context, key string, to any) error

	// MustGet returns data from storage and save it into `to` ignoring errors.
	// Destination argument must be a valid pointer.
	MustGet(ctx context.Context, key string, to any)
}

type fsmContext struct {
	s          Storage
	c          tele.Context
	chat, user int64
}

// NewFSMContext returns new builtin FSM Context.
func NewFSMContext(c tele.Context, storage Storage) Context {
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

func (f *fsmContext) State() (State, error) {
	return f.s.GetState(f.chat, f.user)
}

func (f *fsmContext) Set(state State) error {
	return f.s.SetState(f.chat, f.user, state)
}

func (f *fsmContext) Finish(deleteData bool) error {
	return f.s.ResetState(f.chat, f.user, deleteData)
}

func (f *fsmContext) Update(key string, data any) error {
	return f.s.UpdateData(f.chat, f.user, key, data)
}

func (f *fsmContext) Get(key string, to any) error {
	return f.s.GetData(f.chat, f.user, key, to)
}

func (f *fsmContext) MustGet(key string, to any) {
	_ = f.s.GetData(f.chat, f.user, key, to)
}
