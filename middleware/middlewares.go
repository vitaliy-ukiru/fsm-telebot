// Package middleware is simple middlewares for telebot.
package middleware

import (
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

// ContextKey is key for telebot.Context storage what uses in middleware.
const ContextKey = "fsm"

// FSMContextMiddleware save FSM FSMContext in telebot.Context.
// Recommend use without manager.
func FSMContextMiddleware(storage fsm.Storage) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(ContextKey, fsm.NewFSMContext(c, storage))
			return next(c)
		}
	}
}

// StateFilterMiddleware is filter base on states. Recommended uses only in groups.
// It can be uses if you want handle many endpoints for one state
func StateFilterMiddleware(storage fsm.Storage, want fsm.State) tele.MiddlewareFunc {
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
