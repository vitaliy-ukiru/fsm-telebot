package fsm

import tele "gopkg.in/telebot.v3"

// Group is copy of tele.Group but with support FSM handlers.
// If you use tele.Group then handlers will override by telebot.
// And you can get bugs.
type Group struct {
	m           *Manager
	middlewares []tele.MiddlewareFunc
}

func (g *Group) Use(middlewares ...tele.MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) with(m []tele.MiddlewareFunc) []tele.MiddlewareFunc {
	return append(g.middlewares, m...)
}

func (g *Group) Bind(end interface{}, state State, h Handler, middlewares ...tele.MiddlewareFunc) {
	g.m.Bind(end, state, h, g.with(middlewares)...)
}

func (g *Group) Handle(f Filter, h Handler, middlewares ...tele.MiddlewareFunc) {
	g.m.Handle(f, h, g.with(middlewares)...)
}
