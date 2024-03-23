package fsm

import (
	"context"

	"github.com/vitaliy-ukiru/fsm-telebot/v2/internal/container"
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
	tele "gopkg.in/telebot.v3"
)

type StateMatcher interface {
	MatchState(state State) bool
}

type StateFilter func(state State) bool

func (s StateFilter) MatchState(state State) bool {
	return s(state)
}

func NewMultiStateFilter(states ...State) StateFilter {
	set := container.HashSetFromSlice(states)
	return func(state State) bool {
		return set.Has(state) || set.Has(AnyState)
	}
}

func (m *Manager) Filter(filter StateMatcher) tf.Filter {
	return func(c tele.Context) bool {
		return m.runFilter(c, filter)
	}
}

func (m *Manager) runFilter(c tele.Context, filter StateMatcher) bool {
	fsmCtx := m.mustGetContext(c)

	state, err := fsmCtx.State(context.Background())
	if err != nil {
		c.Bot().OnError(err, c)
		return false
	}

	return filter.MatchState(state)
}
