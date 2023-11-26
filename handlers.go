package fsm

import (
	"context"
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	"github.com/vitaliy-ukiru/fsm-telebot/internal/container"
	tele "gopkg.in/telebot.v3"
)

// handlerMapping contains handlers group separated by endpoint.
type handlerMapping map[string]*handlerList

type handlerList = container.List[handlerEntry]

// handlerEntry representation handler with states, needed for add endpoints correct
// Because telebot uses rule: 1 endpoint = 1 handler.
// But for 1 endpoint allowed more states in our case.
//
// We can use switch-case in handler for check states, but I think not best practice.
type handlerEntry struct {
	states  container.Set[State]
	handler tele.HandlerFunc
}

// add handler to storage, just shortcut.
func (hm handlerMapping) add(endpoint string, h tele.HandlerFunc, states []State) {
	statesSet := container.HashSetFromSlice(states)
	hm.insert(endpoint, handlerEntry{states: statesSet, handler: h})
}

func (hm handlerMapping) insert(endpoint string, entry handlerEntry) {
	if hm[endpoint] == nil {
		hm[endpoint] = new(handlerList)
	}

	hm[endpoint].Insert(entry)
}

// forEndpoint returns handler what filters queries and execute correct handler.
func (m *Manager) forEndpoint(endpoint string) tele.HandlerFunc {
	return func(teleCtx tele.Context) error {
		fsmCtx := m.newContext(teleCtx)

		state, err := fsmCtx.State(context.TODO())
		if err != nil {
			return &ErrHandlerState{Handler: endpoint, Err: err}
		}

		h, ok := m.handlers.find(endpoint, state)
		if !ok {
			return nil
		}

		// middlewares must be executed inside
		// this handler for right work.
		return h.handler(&wrapperContext{teleCtx, fsmCtx})
	}
}

func (hm handlerMapping) find(endpoint string, state State) (handlerEntry, bool) {
	l := hm[endpoint]

	for e := l.Front(); e != nil; e = e.Next() {
		h := e.Value

		if h.states.Has(state) || h.states.Has(AnyState) {
			return h, true
		}
	}

	return handlerEntry{}, false
}

// ErrHandlerState indicates what manager gets error while tired
// get user state in handler.
type ErrHandlerState struct {
	// Handler is the endpoint of the handler
	// where the error occurred.
	Handler string

	// Error what occurred.
	Err error
}

func (e ErrHandlerState) Unwrap() error { return e.Err }

func (e ErrHandlerState) Error() string {
	return fmt.Sprintf(
		"fsm-telebot: get state at handler %s: %v",
		internal.EndpointFormat(e.Handler),
		e.Err,
	)
}
