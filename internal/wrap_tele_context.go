package internal

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

// FsmInternalKey needed for catch
// fsm context requests.
// NOTE: may change to "fsm" for link with middleware/FSMContextMiddleware
const FsmInternalKey = "__fsm"

// WrapperContext wraps telebot context and adds fsm
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
type WrapperContext struct {
	tele.Context
	fsmCtx fsm.Context
}

func NewWrapperContext(context tele.Context, fCtx fsm.Context) *WrapperContext {
	return &WrapperContext{Context: context, fsmCtx: fCtx}
}

func (w *WrapperContext) Get(key string) any {
	if key == FsmInternalKey {
		return w.fsmCtx
	}
	return w.Context.Get(key)
}

func (w *WrapperContext) FSMContext() fsm.Context { return w.fsmCtx }

// TryUnwrapContext tries get fsm.Context from telebot.Context.
func TryUnwrapContext(c tele.Context) (fsm.Context, bool) {
	wrapped, ok := c.(*WrapperContext)
	if ok {
		return wrapped.fsmCtx, true
	}

	ctx, ok := c.Get(FsmInternalKey).(fsm.Context)
	return ctx, ok
}
