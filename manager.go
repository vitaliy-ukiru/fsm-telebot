package fsm

import (
	"github.com/pkg/errors"
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
	bot      *tele.Bot
	group    *tele.Group // handlers will add to group
	store    Storage
	handlers handlerStorage
}

// NewManager returns new Manger.
func NewManager(b *tele.Bot, g *tele.Group, s Storage) *Manager {
	if g == nil {
		g = b.Group()
	}
	return &Manager{
		bot:      b,
		group:    g,
		store:    s,
		handlers: make(handlerStorage),
	}
}

// Group handlers for manger.
func (m *Manager) Group() *tele.Group {
	return m.group
}

func (m *Manager) With(g *tele.Group) *Manager {
	return &Manager{
		bot:      m.bot,
		group:    g,
		store:    m.store,
		handlers: m.handlers,
	}
}

func (m *Manager) NewGroup() *Manager {
	return m.With(m.bot.Group())
}

// Use add middlewares to group.
func (m *Manager) Use(middlewares ...tele.MiddlewareFunc) {
	m.group.Use(middlewares...)
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
	m.handlers.add(endpoint, h, f.States)
	m.group.Handle(endpoint, m.HandlerAdapter(m.handlers.getHandler(endpoint)), middlewares...)

}

// HandlerAdapter create telebot.HandlerFunc object for Handler with FSM FSMContext.
func (m *Manager) HandlerAdapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		return handler(c, NewFSMContext(c, m.s))
	}
}

// Storage returns manger storage instance.
func (m *Manager) Storage() Storage {
	return m.store
}

// GetState returns state for given user in given chat.
func (m *Manager) GetState(chat, user int64) (State, error) {
	return m.store.GetState(chat, user)
}

// SetState sets state for given user in given chat.
func (m *Manager) SetState(chat, user int64, state State) error {
	return m.store.SetState(chat, user, state)
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
		state, err := fsm.State()
		if err != nil {
			return errors.Wrapf(err, "fsm-telebot: get state for endpoint %s", endpoint)
		}

		for _, group := range m[endpoint] {
			if ContainsState(state, group.states...) {
				return group.handler(c, fsm)
			}
		}
		return nil
	}
}
