package fsm

import (
	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM FSMContext
type Handler func(c tele.Context, state Context) error

// ContextMakerFunc alias for function for create new context.
// You can use custom Context implementation.
type ContextMakerFunc func(ctx tele.Context, storage Storage) Context // TODO: add error to return values

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	bot          *tele.Bot
	group        *tele.Group // handlers will add to group
	store        Storage
	handlers     handlerMapping
	contextMaker ContextMakerFunc
	g            []tele.MiddlewareFunc
}

// NewManager returns new Manger.
func NewManager(
	bot *tele.Bot,
	group *tele.Group,
	storage Storage,
	ctxMaker ContextMakerFunc,
) *Manager {
	if group == nil {
		group = bot.Group()
	}
	if ctxMaker == nil {
		ctxMaker = NewFSMContext
	}
	return &Manager{
		bot:          bot,
		group:        group,
		store:        storage,
		contextMaker: ctxMaker,
		handlers:     make(handlerMapping),
	}
}

// Group handlers for manger.
func (m *Manager) Group() *tele.Group {
	return m.group
}

// With return copy of manager with group.
//
// Deprecated: Incorrect behavior with separated groups.
func (m *Manager) With(g *tele.Group) *Manager {
	return &Manager{
		bot:          m.bot,
		group:        g,
		store:        m.store,
		handlers:     m.handlers,
		contextMaker: m.contextMaker,
		g:            m.g,
	}
}

// SetContextMaker sets new context maker to current Manager instance.
func (m *Manager) SetContextMaker(contextMaker ContextMakerFunc) {
	m.contextMaker = contextMaker
}

// NewGroup returns Manager copy with new tele.Group.
//
// Deprecated: Incorrect behavior with separated groups.
func (m *Manager) NewGroup() *Manager {
	return &Manager{
		bot:          m.bot,
		group:        m.bot.Group(),
		store:        m.store,
		handlers:     m.handlers,
		contextMaker: m.contextMaker,
		g:            m.g,
	}
}

// Child returns copy of manager for independence
// adds middlewares.
//
// NOTE: review name before release.
func (m *Manager) Child() *Manager {
	g := make([]tele.MiddlewareFunc, len(m.g))
	copy(g, m.g)
	return &Manager{
		bot:          m.bot,
		group:        m.group,
		store:        m.store,
		handlers:     m.handlers,
		contextMaker: m.contextMaker,
		g:            g,
	}
}

// Use add middlewares to group.
func (m *Manager) Use(middlewares ...tele.MiddlewareFunc) {
	m.g = append(m.g, middlewares...)
}

// Bind adds handler (with FSMContext) with filter on state.
//
// Difference between Bind and Handle methods what Handle require Filter objects.
// And this method can work with only one state.
// If you bind some states see docs to Handle.
func (m *Manager) Bind(end any, state State, h Handler, middlewares ...tele.MiddlewareFunc) {
	m.handle(end, []State{state}, h, middlewares)
}

// Handle adds handler to group chain with filter on states.
// Allowed use more handler for one endpoint.
// If you pass empty slice of states it converters to DefaultState
// Binding some states to one handler
//
//	var ( // types of variables
//		endpoint any // string | tele.CallbackEndpoint
//		states []State
//		handlerFunc fsm.Handler
//	)
//	manager.Handle(fsm.F(endpoint, states...), handlerFunc)
//	// or
//	manager.Handle(fsm.Filter{endpoint, states}, handlerFunc)
func (m *Manager) Handle(f Filter, h Handler, middlewares ...tele.MiddlewareFunc) {
	if len(f.States) == 0 {
		f.States = []State{DefaultState}
	}

	m.handle(f.Endpoint, f.States, h, middlewares)
}

func (m *Manager) handle(
	end any,
	states []State,
	h Handler,
	ms []tele.MiddlewareFunc,
) {
	endpoint := getEndpoint(end)

	// we handles multi handlers in telebot,
	// so need to use middleware here
	wrappedHandler := m.withMiddleware(m.adapter(h), ms)
	m.handlers.add(endpoint, wrappedHandler, states)

	m.group.Handle(
		endpoint,
		m.forEndpoint(endpoint),
	)
}

// withMiddleware join handler middlewares with group middlewares.
func (m *Manager) withMiddleware(h tele.HandlerFunc, ms []tele.MiddlewareFunc) tele.HandlerFunc {
	ms = append(m.g, ms...)

	// I didn’t understand why ApplyMiddleware is called
	// inside the handler, just copied from telebot code.
	return func(c tele.Context) error {
		return internal.ApplyMiddleware(h, ms)(c)
	}
}

// HandlerAdapter create telebot.HandlerFunc object for Handler with FSM FSMContext.
func (m *Manager) HandlerAdapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		return handler(c, m.contextMaker(c, m.store))
	}
}

// adapter wraps internal Handler to telebot.
// difference between HandlerAdapter in support
// wrap context.
func (m *Manager) adapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		fsmCtx, ok := tryUnwrapContext(c)
		if ok {
			return handler(c, fsmCtx)
		}

		// bad case, creating new context
		return handler(c, m.contextMaker(c, m.store))
	}
}

// NewContext creates new FSM Context.
//
// It calls provided ContextMakerFunc.
func (m *Manager) NewContext(teleCtx tele.Context) Context {
	return m.contextMaker(teleCtx, m.store)
}

// Storage returns manger storage instance.
func (m *Manager) Storage() Storage {
	return m.store
}
