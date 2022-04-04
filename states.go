package fsm_telebot

// State objects just string for identification.
// Default state is empty string.
// If state is "*" it corresponds to any state.
type State string

const (
	DefaultState State = ""
	AnyState     State = "*"
)

func (s State) String() string {
	switch s {
	case DefaultState:
		return "State(nil)"
	case AnyState:
		return "State(any)"
	default:
		return string("State(" + s + ")")
	}
}

// Is indicates what state corresponds for other state.
func Is(s State, other State) bool {
	// if current or other state is * => every state equal
	return s == AnyState || other == AnyState || s == other
}

// ContainsState indicates what state contains in given states.
func ContainsState(s State, other ...State) bool {
	for _, state := range other {
		if Is(s, state) {
			return true
		}
	}
	return false
}
