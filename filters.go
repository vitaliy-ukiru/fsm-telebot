package fsm

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

// Filter object. Needs for graceful works with state filters.
type Filter struct {
	Endpoint any
	States   []State
}

// F returns new Filter object.
func F(endpoint any, states ...State) Filter {
	return Filter{Endpoint: endpoint, States: states}
}

// TelebotHandlerForState creates tele.Handler with local filter for given state.
func (m *Manager) TelebotHandlerForState(want State, handler Handler) tele.HandlerFunc {
	return m.HandlerAdapter(func(c tele.Context, state Context) error {
		s, err := state.State()
		if err != nil {
			return &ErrHandlerState{Handler: "Manager.ForState", Err: err}
		}
		if Is(s, want) {
			return handler(c, state)
		}
		return nil
	})
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

// MatchAnyState matches any state.
func MatchAnyState(_ State) bool {
	return true
}

func PrefixMatcher(prefix string) StateMatchFunc {
	return func(state State) bool {
		return strings.HasPrefix(string(state), prefix)
	}
}

func (f Filter) CallbackUnique() string {
	return getEndpoint(f.Endpoint)
}

func getEndpoint(e any) string {
	switch end := e.(type) {
	case string:
		return end
	case tele.CallbackEndpoint:
		return end.CallbackUnique()
	default:
		panic("fsm: telebot: unsupported endpoint")
	}
}
