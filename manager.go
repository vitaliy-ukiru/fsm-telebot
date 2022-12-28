package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM FSMContext
type Handler func(c tele.Context, state Context) error

// fsmHandler representation handler with states, needed for add endpoints correct
// Because telebot uses rule: 1 endpoint = 1 handler. But for 1 endpoint allowed more states.
// We can use switch-case in handler for check states, but I think not best practice.
type fsmHandler struct {
	states  []State
	handler Handler
}

// handlerStorage contains handlers group separated by endpoint.
type handlerStorage map[string][]fsmHandler

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	b *tele.Bot
	g *tele.Group
	s Storage
	h handlerStorage
}

// NewManager returns new Manger.
func NewManager(b *tele.Bot, g *tele.Group, s Storage) *Manager {
	if g == nil {
		g = b.Group()
	}
	return &Manager{b: b, g: g, s: s, h: make(handlerStorage)}
}

// Group handlers for manger.
func (m *Manager) Group() *tele.Group {
	return m.g
}

func (m *Manager) With(g *tele.Group) *Manager {
	return NewManager(m.b, g, m.s)
}

func (m *Manager) NewGroup() *Manager {
	return m.With(m.b.Group())
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
// Allowed use more handler for one endpoint.
func (m *Manager) Handle(f Filter, h Handler, middlewares ...tele.MiddlewareFunc) {
	endpoint := f.CallbackUnique()
	m.h.add(endpoint, h, f.States)
	m.g.Handle(endpoint, m.HandlerAdapter(m.h.getHandler(endpoint)), middlewares...)

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

// GetState returns state for given user in given chat.
func (m *Manager) GetState(chat, user int64) State {
	return m.s.GetState(chat, user)
}

// SetState sets state for given user in given chat.
func (m *Manager) SetState(chat, user int64, state State) {
	m.s.SetState(chat, user, state)
}

// add handler to storage, just shortcut.
func (m handlerStorage) add(endpoint string, h Handler, states []State) {
	m[endpoint] = append(m[endpoint], fsmHandler{
		states:  states,
		handler: h,
	})

}

// getHandler returns handler what filters queries and execute correct handler.
func (m handlerStorage) getHandler(endpoint string) Handler {
	return func(c tele.Context, fsm Context) error {
		state := fsm.State()
		for _, group := range m[endpoint] {
			if ContainsState(state, group.states...) {
				return group.handler(c, fsm)
			}
		}
		return nil
	}
}
