package fsmctx

import (
	"context"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v3"
)

// WithoutContext provide old [fsm.Context] API without context.Context argument.
type WithoutContext struct {
	fsm.Context
}

func NewWithoutContext(context fsm.Context) *WithoutContext {
	return &WithoutContext{Context: context}
}

func (w *WithoutContext) State() (state fsm.State, err error) {
	return w.Context.State(context.Background())
}

func (w *WithoutContext) SetState(state fsm.State) error {
	return w.Context.SetState(context.Background(), state)
}

func (w *WithoutContext) Finish(deleteData bool) error {
	return w.Context.Finish(context.Background(), deleteData)
}

func (w *WithoutContext) Update(key string, data any) error {
	return w.Context.Update(context.Background(), key, data)
}

func (w *WithoutContext) Data(key string, to any) error {
	return w.Context.Data(context.Background(), key, to)
}

func (w *WithoutContext) MustGet(key string, to any) {
	w.Context.MustGet(context.Background(), key, to)
}

type Handler func(c tele.Context, state *WithoutContext) error

func Wrap(h Handler) fsm.Handler {
	return func(teleCtx tele.Context, fsmCtx fsm.Context) error {
		return h(teleCtx, &WithoutContext{Context: fsmCtx})
	}
}
