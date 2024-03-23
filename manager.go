package fsm

import (
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM context.
type Handler func(c tele.Context, state Context) error

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

func (m *Manager) mustGetContext(c tele.Context) Context {
	fsmCtx, ok := tryUnwrapContext(c)
	if ok {
		return fsmCtx
	}
	return m.NewContext(c)
}

func (m *Manager) Storage() Storage {
	return m.store
}

func (m *Manager) Adapt(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		return m.runHandler(c, handler)
	}
}

// HandlerConfig is description of FSM handler.
type HandlerConfig struct {
	Endpoint    any
	OnState     StateMatcher
	Filters     []tf.Filter
	Handler     Handler
	Middlewares []tele.MiddlewareFunc
}

// ---- handler section ----

type HandlerOptionFunc func(hc *HandlerConfig)

type Dispatcher interface {
	Dispatch(tf.Route)
}

// Bind builds handler and to dispatcher. For builtin option see fsmopt pkg.
func (m *Manager) Bind(dp Dispatcher, opts ...HandlerOptionFunc) {
	dp.Dispatch(m.New(opts...))
}

// Handle using telebot-like parameters for adding new handler.
// But it don't supports filters.
func (m *Manager) Handle(
	dp Dispatcher,
	endpoint any,
	onState StateMatcher,
	fn Handler,
	mw ...tele.MiddlewareFunc,
) {
	entity := handlerEntity{
		onState: onState,
		handler: fn,
	}

	route := m.newRoute(endpoint, entity, mw)
	dp.Dispatch(route)
}

func (m *Manager) New(opts ...HandlerOptionFunc) tf.Route {
	hc := new(HandlerConfig)
	for _, opt := range opts {
		opt(hc)
	}

	entity := handlerEntity{
		onState: hc.OnState,
		filters: hc.Filters,
		handler: hc.Handler,
	}
	return m.newRoute(hc.Endpoint, entity, hc.Middlewares)
}

func (m *Manager) newRoute(e any, entity handlerEntity, mw []tele.MiddlewareFunc) tf.Route {
	return tf.Route{
		Endpoint: e,
		Handler: &fsmHandler{
			handlerEntity: entity,
			manager:       m,
		},
		Middlewares: mw,
	}
}
