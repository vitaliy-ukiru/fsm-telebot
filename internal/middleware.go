package internal

import tele "gopkg.in/telebot.v3"

// ApplyMiddleware is copy of tele.applyMiddleware.
// For support middlewares in handlers we need packs middlewares independently
// of telebot.
//
// In fact, I'm starting to worry that the package is turning into a
// complete copy (rework) of the telebot. Maybe it's worth stopping?
func ApplyMiddleware(h tele.HandlerFunc, m []tele.MiddlewareFunc) tele.HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// JoinMiddlewares copy two slice to new or return source non-empty slice.
func JoinMiddlewares(a, b []tele.MiddlewareFunc) []tele.MiddlewareFunc {
	aLen := len(a)
	if aLen == 0 {
		return b
	}
	bLen := len(b)
	if bLen == 0 {
		return a
	}

	mw := make([]tele.MiddlewareFunc, aLen+bLen)
	copy(mw, a)
	copy(mw[aLen:], b)
	return mw
}
