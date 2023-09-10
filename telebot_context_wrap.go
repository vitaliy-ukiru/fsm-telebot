package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// fsmInternalKey needed for catch
// fsm context requests.
// NOTE: may change to "fsm" for link with middleware/FSMContextMiddleware
const fsmInternalKey = "__fsm"

// wrapperContext wraps telebot context and adds fsm
// context inside.
// By this wrapper you can get context from any handler
// under this wrapper.
//
// But this very not recommend, because it more internal
// mechanism.
// Also, it will have many consequences, because others
// middlewares can set data with same key (I think it
// would have been done by accident).
//
// The developers of the package make no guarantee
// of use outside of this package.
type wrapperContext[C Context] struct {
	tele.Context
	fsmCtx C
}

func (w *wrapperContext[C]) Get(key string) any {
	if key == fsmInternalKey {
		return w.fsmCtx
	}
	return w.Context.Get(key)
}

func (w *wrapperContext[C]) FSMContext() C { return w.fsmCtx }

// tryUnwrapContext tries get fsm.Context from telebot.Context.
func tryUnwrapContext[C Context](c tele.Context) (C, bool) {
	wrapped, ok := c.(*wrapperContext[C])
	if ok {
		return wrapped.fsmCtx, true
	}

	ctx, ok := c.Get(fsmInternalKey).(C)
	return ctx, ok
}
