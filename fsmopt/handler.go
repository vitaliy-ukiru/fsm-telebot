package fsmopt

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
	tele "gopkg.in/telebot.v3"
)

func Use(mw ...tele.MiddlewareFunc) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.Middlewares = mw
	}
}

func On(e any) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.Endpoint = e
	}
}

func Filter(filters ...tf.Filter) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.Filters = filters
	}
}

func Do(h fsm.Handler) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.Handler = h
	}
}

func OnStates(states ...fsm.State) fsm.HandlerOption {
	var filter fsm.StateMatcher
	switch len(states) {
	case 0:
		filter = fsm.DefaultState
	case 1:
		filter = states[0]
	default:
		filter = fsm.NewMultiStateFilter(states...)
	}
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = filter
	}
}

func FilterState(filter fsm.StateFilter) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = filter
	}
}

func MatchState(matcher fsm.StateMatcher) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = matcher
	}
}

func Config(config fsm.HandlerConfig) fsm.HandlerOption {
	return func(hc *fsm.HandlerConfig) {
		*hc = config
	}
}
