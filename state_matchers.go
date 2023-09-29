package fsm

import (
	"github.com/vitaliy-ukiru/fsm-telebot/internal"
)

func (s State) MatchState(state State) bool { return Is(s, state) }

func (m StateMatchFunc) MatchState(state State) bool {
	return m(state)
}

type setStatesMatcher struct {
	states internal.HashSet[State]
}

func newStateMatcherSlice(states []State) setStatesMatcher {
	hs := make(internal.HashSet[State])
	for _, state := range states {
		hs.Add(state)
	}
	return setStatesMatcher{states: hs}
}

func (s setStatesMatcher) MatchState(state State) bool {
	return s.states.Has(state) || s.states.Has(AnyState)
}

// Matcher returns new matcher object, what will
// match states from group.
func (s *StateGroup) Matcher() StateMatcher {
	return newStateMatcherSlice(s.States)
}

// MatchStates returns matcher based on hashset. It
// also will check match on [AnyState].
func MatchStates(states ...State) StateMatcher {
	return newStateMatcherSlice(states)
}
