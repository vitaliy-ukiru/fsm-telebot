package fsm

import (
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
	tele "gopkg.in/telebot.v3"
)

func (m *Manager) runHandler(c tele.Context, handler Handler) error {
	fsmCtx, ok := m.mustGetContext(c)
	// don't run handler if can't get context
	if !ok || fsmCtx == nil {
		return nil
	}
	return handler(c, fsmCtx)
}

// WrapContext is middleware for wrapping fsm context. It helps to create
// context only one time for update and make small allocation optimization.
// FSM will unwrap this context in internal mechanic.
func (m *Manager) WrapContext(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := newWrapperContext(c, m.NewContext(c))
		return next(ctx)
	}
}

type fsmHandler struct {
	onState StateMatcher
	filter  tf.Filter
	handler Handler
	manager *Manager
}

func (fh fsmHandler) Check(c tele.Context) bool {
	// skip state filter on nil
	if fh.onState != nil && !fh.manager.runFilter(c, fh.onState) {
		return false
	}

	return fh.filter == nil || fh.filter(c)
}

func (fh fsmHandler) Execute(c tele.Context) error {
	return fh.manager.runHandler(c, fh.handler)
}
