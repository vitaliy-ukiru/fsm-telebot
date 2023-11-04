package ext

import (
	"errors"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

func FSM(handler fsm.Handler) Handler {
	return func(c tele.Context, ext Context) error {
		return handler(c, fsm.Context(ext))
	}
}

type StateHandler func(c tele.Context, ext Context) (fsm.State, error)

var SkipState = errors.New("step: skip state step")

func WithState(fn StateHandler) Handler {
	return func(c tele.Context, ext Context) error {
		state, err := fn(c, ext)
		if err != nil {
			if errors.Is(err, SkipState) {
				return nil
			}
			return err
		}
		return ext.Set(state)
	}
}

func NewStateStep(endpoint any, state fsm.State, handler StateHandler) Step {
	return Step{endpoint, state, WithState(handler)}
}
