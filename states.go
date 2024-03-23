package fsm

import (
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot/v2/internal/container"
)

// State objects just string for identification.
//
// Default state is empty string.
// If state is "*" it corresponds to any state.
type State string

const (
	DefaultState State = ""
	AnyState     State = "*"
)

func (s State) MatchState(other State) bool {
	return s == other || (s == AnyState || other == AnyState)
}

func (s State) GoString() string {
	switch s {
	case DefaultState:
		return "State(default)"
	case AnyState:
		return "State(any)"
	default:
		return string("State(" + s + ")")
	}
}

// StateGroup storages states with custom prefix.
//
// It can use in filter like and handled via Manager.Handle:
//
//	group := fsm.NewStateGroup("adm", "State0", "State1")
//	filter := fsm.F("/cmd", group.States...)
type StateGroup struct {
	prefix string
	group  *container.LinkedHashSet[State]
}

func (sg StateGroup) Prefix() string {
	return sg.prefix
}

// NewStateGroup returns new StateGroup.
func NewStateGroup(prefix string, states ...State) *StateGroup {
	sgPrefix := State(prefix + stateGroupSep)
	for i := 0; i < len(states); i++ {
		states[i] = sgPrefix + states[i]
	}
	return &StateGroup{
		prefix: prefix,
		group:  container.NewLinkedHashSet(states...),
	}
}

const stateGroupSep = ":"

func (sg StateGroup) New(name string) State {
	prefix := sg.prefix + stateGroupSep
	if !strings.HasPrefix(name, prefix) {
		name = prefix + name
	}
	state := State(name)
	sg.group.Add(state)
	return state
}

func (sg StateGroup) Next(s State) State {
	node := sg.group.Item(s)
	if node == nil {
		return DefaultState
	}

	next := node.Next()
	if next == nil {
		return DefaultState
	}

	return next.Value()
}

func (sg StateGroup) Prev(s State) State {
	node := sg.group.Item(s)
	if node == nil {
		return DefaultState
	}

	prev := node.Prev()
	if prev == nil {
		return DefaultState
	}

	return prev.Value()
}

func (sg StateGroup) MatchState(state State) bool {
	return sg.group.Has(state)
}
