package fsmopt

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
	tele "gopkg.in/telebot.v3"
)

func Use(mw ...tele.MiddlewareFunc) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.Middlewares = mw
	}
}

func On(e any) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.Endpoint = e
	}
}

func Filter(filters ...tf.Filter) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.Filters = filters
	}
}

func Do(h fsm.Handler) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.Handler = h
	}
}

func OnStates(states ...fsm.State) fsm.HandlerOptionFunc {
	var filter fsm.StateFilter
	switch len(states) {
	case 0:
		filter = fsm.NewSingleStateFilter(fsm.DefaultState)
	case 1:
		filter = fsm.NewSingleStateFilter(states[0])
	default:
		filter = fsm.NewMultiStateFilter(states...)
	}
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = filter
	}
}

func FilterState(filter fsm.StateFilter) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = filter
	}
}

func MatchState(matcher fsm.StateMatcher) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		hc.OnState = matcher
	}
}

func Config(config fsm.HandlerConfig) fsm.HandlerOptionFunc {
	return func(hc *fsm.HandlerConfig) {
		*hc = config
	}
}
