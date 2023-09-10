package middleware

import (
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

// ContextKey is key for telebot.Context storage what uses in middleware.
var ContextKey = "fsm"

// FSMContextMiddleware save FSM context in telebot.Context.
// Recommend use without manager.
func FSMContextMiddleware(storage fsm.Storage) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(ContextKey, fsm.NewFSMContext(c, storage))
			return next(c)
		}
	}
}
