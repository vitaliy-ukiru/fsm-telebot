package fsm

import (
	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM context.
type Handler[C Context] func(c tele.Context, state C) error

// ContextMakerFunc alias for function for create new context.
// You can use custom Context implementation.
type ContextMakerFunc[C Context] func(ctx tele.Context, storage Storage) C // TODO: add error to return values

// Manager is object for managing FSM, binding handlers.
type Manager[C Context, S Storage] struct {
	bot          *tele.Bot
	group        *tele.Group // handlers will add to group
	store        S
	handlers     handlerMapping
	contextMaker ContextMakerFunc[C]
	g            []tele.MiddlewareFunc
}

// NewManager returns new Manger.
func NewManager[C Context, S Storage](
	bot *tele.Bot,
	group *tele.Group,
	storage S,
	ctxMaker ContextMakerFunc[C],
) *Manager[C, S] {
	if group == nil {
		group = bot.Group()
	}
	return &Manager[C, S]{
		bot:          bot,
		group:        group,
		store:        storage,
		contextMaker: ctxMaker,
		handlers:     make(handlerMapping),
	}
}

// Group handlers for manger.
func (m *Manager[C, S]) Group() *tele.Group {
	return m.group
}

// With return copy of manager with group.
//
// Deprecated: Incorrect behavior with separated groups.
func (m *Manager[C, S]) With(g *tele.Group) *Manager[C, S] {
	newM := *m
	newM.group = g
	return &newM
}

// SetContextMaker sets new context maker to current Manager instance.
func (m *Manager[C, S]) SetContextMaker(contextMaker ContextMakerFunc[C]) {
	m.contextMaker = contextMaker
}

// NewGroup returns manager child with copy
// of middleware group. Adding middlewares in
// new group doesn't affect the parent.
func (m *Manager[C, S]) NewGroup() *Manager[C, S] {
	newM := *m
	newM.g = make([]tele.MiddlewareFunc, len(m.g))
	copy(newM.g, m.g)
	return &newM
}

// Use add middlewares to group.
func (m *Manager[C, S]) Use(middlewares ...tele.MiddlewareFunc) {
	m.g = append(m.g, middlewares...)
}

// Bind adds handler (with FSM context argument) with filter on state.
//
// Difference between Bind and Handle methods what Handle require Filter objects.
// And this method can work with only one state.
// If you bind some states see docs to Handle.
func (m *Manager[C, S]) Bind(end any, state State, h Handler[C], middlewares ...tele.MiddlewareFunc) {
	m.handle(end, []State{state}, h, middlewares)
}

// Handle adds handler to group chain with filter on states.
// Allowed use more handler for one endpoint.
// If you pass empty slice of states it converters to DefaultState
// Binding some states to one handler.
//
//	var ( // types of variables
//		endpoint any // string | tele.CallbackEndpoint
//		states []State
//		handlerFunc fsm.Handler
//	)
//	manager.Handle(fsm.F(endpoint, states...), handlerFunc)
//	// or
//	manager.Handle(fsm.Filter{endpoint, states}, handlerFunc)
func (m *Manager[C, S]) Handle(f Filter, h Handler[C], middlewares ...tele.MiddlewareFunc) {
	if len(f.States) == 0 {
		f.States = []State{DefaultState}
	}

	m.handle(f.Endpoint, f.States, h, middlewares)
}

func (m *Manager[C, S]) handle(
	end any,
	states []State,
	h Handler[C],
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
func (m *Manager[C, S]) withMiddleware(h tele.HandlerFunc, ms []tele.MiddlewareFunc) tele.HandlerFunc {
	ms = append(m.g, ms...)

	// I didn’t understand why ApplyMiddleware is called
	// inside the handler, just copied from telebot code.
	return func(c tele.Context) error {
		return internal.ApplyMiddleware(h, ms)(c)
	}
}

// HandlerAdapter create telebot.HandlerFunc object
// for Handler with FSM context.
//
// Used for external purposes only outside handlers chain.
// Example: access to context without manager handlers.
// Use only as directed and if you know what you are doing.
func (m *Manager[C, S]) HandlerAdapter(handler Handler[C]) tele.HandlerFunc {
	return func(c tele.Context) error {
		return handler(c, m.contextMaker(c, m.store))
	}
}

// adapter wraps internal Handler to telebot.
// difference between HandlerAdapter in support
// wrap context.
func (m *Manager[C, S]) adapter(handler Handler[C]) tele.HandlerFunc {
	return func(c tele.Context) error {
		fsmCtx, ok := tryUnwrapContext[C](c)
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
func (m *Manager[C, S]) NewContext(teleCtx tele.Context) C {
	return m.contextMaker(teleCtx, m.store)
}

// Storage returns manger storage instance.
func (m *Manager[C, S]) Storage() S {
	return m.store
}
