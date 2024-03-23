package fsm

import (
	"context"
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
}

type fsmContext struct {
	storage Storage
	key     StorageKey
}

// NewFSMContext returns new builtin FSM Context.
func NewFSMContext(storage Storage, key StorageKey) Context {
	return &fsmContext{
		storage: storage,
		key:     key,
	}
}

func (f *fsmContext) State(ctx context.Context) (State, error) {
	return f.storage.GetState(ctx, f.key)
}

func (f *fsmContext) SetState(ctx context.Context, state State) error {
	return f.storage.SetState(ctx, f.key, state)
}

func (f *fsmContext) Finish(ctx context.Context, deleteData bool) error {
	return f.storage.ResetState(ctx, f.key, deleteData)
}

func (f *fsmContext) Update(ctx context.Context, key string, data any) error {
	return f.storage.UpdateData(ctx, f.key, key, data)
}

func (f *fsmContext) Data(ctx context.Context, key string, to any) error {
	return f.storage.GetData(ctx, f.key, key, to)
}
