package fsm

import (
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

// handlerStorage contains handlers group separated by endpoint.
type handlerStorage map[string][]handlerEntry

// handlerEntry representation handler with states, needed for add endpoints correct
// Because telebot uses rule: 1 endpoint = 1 handler. But for 1 endpoint allowed more states.
// We can use switch-case in handler for check states, but I think not best practice.
type handlerEntry struct {
	states  map[State]struct{}
	handler Handler
}

func (h handlerEntry) match(state State) bool {
	_, ok := h.states[state]
	return ok
}

// add handler to storage, just shortcut.
func (m handlerStorage) add(endpoint string, h Handler, states []State) {
	statesSet := make(map[State]struct{})
	for _, state := range states {
		statesSet[state] = struct{}{}
	}

	m[endpoint] = append(m[endpoint], handlerEntry{
		states:  statesSet,
		handler: h,
	})
}

// forEndpoint returns handler what filters queries and execute correct handler.
func (m handlerStorage) forEndpoint(endpoint string) Handler {
	return func(teleCtx telebot.Context, fsmCtx Context) error {
		state, err := fsmCtx.State()
		if err != nil {
			return errors.Wrapf(err, "fsm-telebot: get state for endpoint %s", endpoint)
		}

		for _, group := range m[endpoint] {
			if group.match(state) || group.match(AnyState) {
				return group.handler(teleCtx, fsmCtx)
			}
		}
		return nil
	}
}
