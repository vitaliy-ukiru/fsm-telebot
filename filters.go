package fsm

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

// Filter object. Needs for graceful works with state filters.
type Filter struct {
	Endpoint interface{}
	States   []State
}

// F returns new Filter object.
func F(endpoint interface{}, states ...State) Filter {
	if len(states) == 0 {
		states = []State{DefaultState}
	}
	return Filter{Endpoint: endpoint, States: states}
}

// TelebotHandlerForState creates tele.Handler with local filter for given state.
func (m *Manager) TelebotHandlerForState(want State, handler Handler) tele.HandlerFunc {
	return m.TelebotHandlerForStates(handler, want)
}

// TelebotHandlerForStates creates a handler with local filter
// for current state to check for presence in given states.
func (m *Manager) TelebotHandlerForStates(h Handler, states ...State) tele.HandlerFunc {
	return m.HandlerAdapter(func(c tele.Context, state Context) error {
		s, err := state.State()
		if err != nil {
			return &ErrHandlerState{Handler: "Manager.ForStates", Err: err}
		}
		if ContainsState(s, states...) {
			return h(c, state)
		}
		return nil
	})
}

func (f Filter) CallbackUnique() string {
	return callbackUnique(f.Endpoint)
}

func callbackUnique(endpoint interface{}) string {
	switch end := endpoint.(type) {
	case string:
		return end
	case tele.CallbackEndpoint:
		return end.CallbackUnique()
	default:
		panic("fsm: telebot: unsupported endpoint")
	}
}

// PrefixFilter filters state on have prefix
// with separated symbol.
func PrefixFilter(prefix string) StateMatchFunc {
	prefix = prefix + "@"
	return func(state State) bool {
		return strings.HasPrefix(string(state), prefix)
	}
}

// MatchAnyState matches any state.
func MatchAnyState(_ State) bool {
	return true
}
