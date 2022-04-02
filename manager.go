package fsm_telebot

import (
	"github.com/vitaliy-ukiru/fsm-telebot/middleware"
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM FSMContext
type Handler func(c tele.Context, state FSMContext) error

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	g *tele.Group
	s Storage
}

// NewManager returns new Manger
func NewManager(group *tele.Group, storage Storage) *Manager {
	return &Manager{g: group, s: storage}
}

// Group handlers for manger.
func (m *Manager) Group() *tele.Group {
	return m.g
}

// Use add middlewares to group.
func (m *Manager) Use(middlewares ...tele.MiddlewareFunc) {
	m.g.Use(middlewares...)
}

// Handle binding handler (with FSMContext) with filter on state.
func (m *Manager) Handle(end interface{}, state State, h Handler, middlewares ...tele.MiddlewareFunc) {
	m.g.Handle(end, m.StateFilter(state, h), middlewares...)

}

// HandlerAdapter create telebot.HandlerFunc object for Handler with FSM FSMContext.
func (m *Manager) HandlerAdapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		fsmCtx, ok := c.Get(middleware.ContextKey).(FSMContext)
		if fsmCtx == nil || !ok {
			fsmCtx = NewFSMContext(c, m.s)
		}
		return handler(c, fsmCtx)
	}
}

// StateFilter filtering updates with current state and execute handler.
func (m *Manager) StateFilter(state State, handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		if m.GetState(c.Chat().ID, c.Sender().ID).Is(state) {
			return handler(c, NewFSMContext(c, m.s))
		}
		return nil
	}
}

// Storage returns manger storage instance.
func (m *Manager) Storage() Storage {
	return m.s
}

// GetState returns state for current user in current chat.
func (m *Manager) GetState(chat, user int64) State {
	return m.s.GetState(chat, user)
}

// SetState sets state for current user in current chat.
func (m *Manager) SetState(chat, user int64, state State) error {
	return m.s.SetState(chat, user, state)
}
