package fsm_telebot

import (
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

// Bind adds handler (with FSMContext) with filter on state.
//
// Difference between Bind and Handle methods what Handle require Filter objects.
// And this method can work with only one state.
func (m *Manager) Bind(end interface{}, state State, h Handler, middlewares ...tele.MiddlewareFunc) {
	m.Handle(F(end, state), h, middlewares...)

}

// Handle adds handler to group chain with filter on states.
func (m *Manager) Handle(f Filter, h Handler, middlewares ...tele.MiddlewareFunc) {
	m.g.Handle(f.Endpoint, m.ForStates(h, f.States...), middlewares...)
}

// HandlerAdapter create telebot.HandlerFunc object for Handler with FSM FSMContext.
func (m *Manager) HandlerAdapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		return handler(c, NewFSMContext(c, m.s))
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
