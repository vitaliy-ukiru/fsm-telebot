package fsmctx

import (
	"context"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

type StorageKey struct {
	Chat int64
	User int64
}

type CtxStorage interface {
	GetStateCtx(ctx context.Context, target StorageKey) (fsm.State, error)
	SetStateCtx(ctx context.Context, target StorageKey, state fsm.State) error
	ResetStateCtx(ctx context.Context, target StorageKey, withData bool) error
	UpdateDataCtx(ctx context.Context, target StorageKey, key string, data any) error
	GetDataCtx(ctx context.Context, target StorageKey, key string, to any) error
	CloseCtx(ctx context.Context) error
}

type CtxContext interface {
	StateCtx(ctx context.Context) (fsm.State, error)
	SetCtx(ctx context.Context, state fsm.State) error
	// not full example
}

type WithContext struct {
	// for implement default fsm.Context
	*fsm.BuiltinContext

	k StorageKey
	s CtxStorage
}

func NewWithContext(c tele.Context, s fsm.Storage) *WithContext {
	return &WithContext{
		BuiltinContext: fsm.NewFSMContext(c, s),
		k:              StorageKey{Chat: c.Chat().ID, User: c.Sender().ID},
		s:              s.(CtxStorage),
	}
}

func (w *WithContext) StateCtx(ctx context.Context) (fsm.State, error) {
	return w.s.GetStateCtx(ctx, w.k)
}

func (w *WithContext) SetCtx(ctx context.Context, state fsm.State) error {
	return w.s.SetStateCtx(ctx, w.k, state)
}

type ContextEntry struct {
	ctx    context.Context
	parent *WithContext
}

func (w *WithContext) With(ctx context.Context) *ContextEntry {
	return &ContextEntry{ctx: ctx, parent: w}
}

func (c *ContextEntry) State() (fsm.State, error) {
	return c.parent.StateCtx(c.ctx)
}

func (c *ContextEntry) SetState(state fsm.State) error {
	return c.parent.SetCtx(c.ctx, state)
}
