package fsmctx_test

import (
	"context"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/adapter/fsmctx"
	tele "gopkg.in/telebot.v3"
)

func ExampleWrap() {
	const ContextState fsm.State = "context_state"
	fsmctx.Wrap(func(c tele.Context, stateCtx *fsmctx.WithoutContext) error {
		state, err := stateCtx.State()
		if err != nil {
			return err
		}

		if state == ContextState {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// if you need use ctx arg you can use field stateCtx.WithoutContext
			if err := stateCtx.Context.SetState(ctx, "fast_state"); err != nil {
				return err
			}
		}

		return c.Send("Bye!")
	})
}
