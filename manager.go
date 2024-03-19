package fsm

import (
	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM context.
type Handler func(c tele.Context, state Context) error

// ContextMakerFunc alias for function for create new context.
// ContextFactoryFunc alias for function for create new context.
// You can use custom Context implementation.
type ContextFactoryFunc func(storage Storage, key StorageKey) Context

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	store          Storage
	strategy       Strategy
	contextFactory ContextFactoryFunc
}

func New(store Storage, strategy Strategy, contextMaker ContextFactoryFunc) *Manager {
	if contextMaker == nil {
		contextMaker = NewFSMContext
	}
	return &Manager{
		store:          store,
		strategy:       strategy,
		contextFactory: contextMaker,
	}
}

// NewContext creates new FSM Context.
//
// It calls provided ContextFactoryFunc.
func (m *Manager) NewContext(ctx tele.Context) Context {
	key := ExtractKeyWithStrategy(ctx, m.strategy)
	return m.contextFactory(m.store, key)
}

// Storage returns manger storage instance.
func (m *Manager) Storage() Storage {
	return m.store
}
