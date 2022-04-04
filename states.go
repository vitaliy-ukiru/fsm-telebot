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
func (s State) Is(other State) bool {
	// if current or other state is * => every state equal
	return s == AnyState || other == AnyState || s == other
}
