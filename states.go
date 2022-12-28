package fsm

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
	return s == other || (s == AnyState || other == AnyState)
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

// StateGroup storages states with custom prefix.
//
// It can use in filter like:
// 	group := fsm.NewStateGroup("adm", "State0", "State1")
//	filter := fsm.F("/cmd", group.States...)
type StateGroup struct {
	Prefix string
	States []State
}

// NewStateGroup returns new StateGroup.
func NewStateGroup(prefix string, states ...State) *StateGroup {
	return &StateGroup{
		Prefix: prefix,
		States: states,
	}
}

// New returns new state with group prefix and add to group states.
func (s *StateGroup) New(name string) (state State) {
	state = State(s.Prefix + "@" + name)
	s.States = append(s.States, state)
	return
}
