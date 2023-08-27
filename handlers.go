package fsm

import (
	"container/list"
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	tele "gopkg.in/telebot.v3"
)

// handlerStorage contains handlers group separated by endpoint.
type handlerStorage map[string]*list.List

// handlerEntry representation handler with states, needed for add endpoints correct
// Because telebot uses rule: 1 endpoint = 1 handler.
// But for 1 endpoint allowed more states in our case.
//
// We can use switch-case in handler for check states, but I think not best practice.
type handlerEntry struct {
	states  statesHashset
	matchFn StateMatchFunc
	handler Handler
}

// add handler to storage, just shortcut.
func (m handlerStorage) add(endpoint string, h Handler, states []State) {
	statesSet := newHashsetFromSlice(states)
	m.insert(endpoint, handlerEntry{states: statesSet, handler: h})
}
func (m handlerStorage) addFunc(endpoint string, h Handler, matchFunc StateMatchFunc) {
	m.insert(endpoint, handlerEntry{matchFn: matchFunc, handler: h})
}

func (m handlerStorage) insert(endpoint string, entry handlerEntry) {
	if m[endpoint] == nil {
		m[endpoint] = list.New()
	}

	m[endpoint].PushBack(entry)
}

// forEndpoint returns handler what filters queries and execute correct handler.
func (m handlerStorage) forEndpoint(endpoint string) Handler {
	return func(teleCtx tele.Context, fsmCtx Context) error {
		state, err := fsmCtx.State()
		if err != nil {
			return &ErrHandlerState{Handler: endpoint, Err: err}
		}

		l := m[endpoint]

		for e := l.Front(); e != nil; e = e.Next() {
			h := e.Value.(handlerEntry)

			matchState := (h.matchFn != nil && h.matchFn(state)) ||
				(h.states.Has(state) || h.states.Has(AnyState))

			if matchState {
				return h.handler(teleCtx, fsmCtx)
			}
		}
		return nil
	}
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
