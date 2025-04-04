package fsm

import (
	filterspkg "github.com/vitaliy-ukiru/telebot-filter/v2/pkg/filters"
	tf "github.com/vitaliy-ukiru/telebot-filter/v2/telefilter"
	tele "gopkg.in/telebot.v4"
)

// Handler is function for handling updates
// with access to user state via fsm context.
type Handler func(c tele.Context, state Context) error

// ContextFactory creates new FSM context.
type ContextFactory func(storage Storage, key StorageKey) Context

// StateFilterProcessor must return filter result.
// Processor should handle problems by self or
// gives this task to external systems.
type StateFilterProcessor func(c tele.Context, fsmCtx Context, matcher StateMatcher) bool

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	store           Storage
	strategy        Strategy
	contextFactory  ContextFactory
	filterProcessor StateFilterProcessor
}

// Settings provides configuration for [Manager].
// Works via [ManagerOption]
type Settings struct {
	Strategy
	ContextFactory
	StateFilterProcessor
}

type ManagerOption func(*Settings)

func New(storage Storage, opts ...ManagerOption) *Manager {
	cfg := new(Settings)

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ContextFactory == nil {
		cfg.ContextFactory = NewFSMContext
	}

	if cfg.StateFilterProcessor == nil {
		cfg.StateFilterProcessor = DefaultFilterProcessor
	}

	return &Manager{
		store:           storage,
		strategy:        cfg.Strategy,
		contextFactory:  cfg.ContextFactory,
		filterProcessor: cfg.StateFilterProcessor,
	}
}

// NewContext creates new FSM Context.
// Context will create via call context factory.
//
// If key will be non-present it will return (nil, false)
func (m *Manager) NewContext(ctx tele.Context) (Context, bool) {
	key, ok := extractKeyWithStrategy(ctx, m.strategy)
	if !ok {
		return nil, false
	}

	context := m.contextFactory(m.store, key)
	return context, context != nil
}

func (m *Manager) mustGetContext(c tele.Context) (Context, bool) {
	fsmCtx, ok := tryUnwrapContext(c)
	if ok {
		return fsmCtx, fsmCtx != nil
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

// ---- handler section ----

// HandlerConfig is description of FSM handler.
type HandlerConfig struct {
	Endpoint    any
	OnState     StateMatcher
	Filters     []tf.Filter
	Handler     Handler
	Middlewares []tele.MiddlewareFunc
}

type HandlerOption func(hc *HandlerConfig)

type Dispatcher interface {
	Dispatch(tf.Route)
}

// Bind builds handler and to dispatcher. For builtin option see fsmopt pkg.
func (m *Manager) Bind(dp Dispatcher, opts ...HandlerOption) {
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
	handler := fsmHandler{
		onState: onState,
		handler: fn,
	}

	route := m.newRoute(endpoint, handler, mw)
	dp.Dispatch(route)
}

func (m *Manager) New(opts ...HandlerOption) tf.Route {
	hc := new(HandlerConfig)
	for _, opt := range opts {
		opt(hc)
	}

	handler := fsmHandler{
		onState: hc.OnState,
		filter:  combineFilters(hc.Filters),
		handler: hc.Handler,
	}
	return m.newRoute(hc.Endpoint, handler, hc.Middlewares)
}

func (m *Manager) newRoute(e any, handler fsmHandler, mw []tele.MiddlewareFunc) tf.Route {
	handler.manager = m

	return tf.Route{
		Endpoint:    e,
		Handler:     handler,
		Middlewares: mw,
	}
}

func combineFilters(filters []tf.Filter) tf.Filter {
	n := len(filters)
	switch n {
	case 0:
		return nil
	case 1:
		return filters[0]
	}
	return filterspkg.All(filters...)
}
