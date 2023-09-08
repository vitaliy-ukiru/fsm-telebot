package fsm

import (
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/internal"
	tele "gopkg.in/telebot.v3"
)

// handlerMapping contains handlers group separated by endpoint.
type handlerMapping map[string]*internal.List[handlerEntry]

// handlerEntry representation handler with states, needed for add endpoints correct
// Because telebot uses rule: 1 endpoint = 1 handler.
// But for 1 endpoint allowed more states in our case.
//
// We can use switch-case in handler for check states, but I think not best practice.
type handlerEntry struct {
	states  internal.HashSet[State]
	handler tele.HandlerFunc
}

// add handler to storage, just shortcut.
func (hm handlerMapping) add(endpoint string, h tele.HandlerFunc, states []State) {
	statesSet := internal.HashSetFromSlice(states)
	hm.insert(endpoint, handlerEntry{states: statesSet, handler: h})
}

func (hm handlerMapping) insert(endpoint string, entry handlerEntry) {
	if hm[endpoint] == nil {
		hm[endpoint] = new(internal.List[handlerEntry])
	}

	hm[endpoint].Insert(entry)
}

// forEndpoint returns handler what filters queries and execute correct handler.
func (m handlerStorage) forEndpoint(endpoint string) Handler {
	return func(teleCtx tele.Context, fsmCtx Context) error {
		state, err := fsmCtx.State()
		if err != nil {
			return &ErrHandlerState{Handler: endpoint, Err: err}
		}

		h, ok := m.findHandler(endpoint, state)
		if !ok {

			return nil
		}
		return h.handler(teleCtx, fsmCtx)

	}
}

func (hm handlerMapping) findHandler(endpoint string, state State) (handlerEntry, bool) {
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
