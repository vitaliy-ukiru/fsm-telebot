package middleware

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

// ContextKey is key for telebot.Context storage what uses in middleware.
const ContextKey = "fsm"

// FSMContextMiddleware save FSM FSMContext in telebot.Context.
func FSMContextMiddleware(storage fsm_telebot.Storage) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(ContextKey, fsm_telebot.NewFSMContext(c, storage))
			return next(c)
		}
	}
}

// StateFilterMiddleware is filter base on states. Recommended uses only in groups.
func StateFilterMiddleware(storage fsm_telebot.Storage, want fsm_telebot.State) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			currentState := storage.GetState(c.Chat().ID, c.Sender().ID)
			if currentState.Is(want) {
				return next(c)
			}
			return nil
		}
	}
}
