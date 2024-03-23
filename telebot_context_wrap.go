package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// fsmInternalKey needed for catch context requests.
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
type wrapperContext struct {
	tele.Context
	fsmCtx Context
}

func newWrapperContext(context tele.Context, fsmCtx Context) *wrapperContext {
	return &wrapperContext{Context: context, fsmCtx: fsmCtx}
}

func (w *wrapperContext) Get(key string) any {
	if key == fsmInternalKey {
		return w.fsmCtx
	}
	return w.Context.Get(key)
}

func (w *wrapperContext) FSMContext() Context { return w.fsmCtx }

// tryUnwrapContext tries get fsm.Context from telebot.Context.
func tryUnwrapContext(c tele.Context) (Context, bool) {
	wrapped, ok := c.(*wrapperContext)
	if ok {
		return wrapped.fsmCtx, true
	}

	ctx, ok := c.Get(fsmInternalKey).(Context)
	return ctx, ok
}
