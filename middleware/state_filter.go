// Package middleware is simple middlewares for telebot.
package middleware

import (
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

// StateFilterMiddleware is filter base on states. Recommended uses only in groups.
// It can be uses if you want handle many endpoints for one state
func StateFilterMiddleware(storage fsm.Storage, want fsm.State) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			currentState, err := storage.GetState(c.Chat().ID, c.Sender().ID)
			if err != nil {
				return err
			}
			if fsm.Is(currentState, want) {
				return next(c)
			}
			return nil
		}
	}
}
