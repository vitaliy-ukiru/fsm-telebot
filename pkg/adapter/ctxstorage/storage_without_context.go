package ctxstorage

import (
	"context"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

// StorageWithoutContext is old storage model, that don't support context flow.
//
// Need sync with fsm.Storage interface.
type StorageWithoutContext interface {
	GetState(targetKey fsm.StorageKey) (fsm.State, error)
	SetState(targetKey fsm.StorageKey, state fsm.State) error
	ResetState(targetKey fsm.StorageKey, withData bool) error
	UpdateData(targetKey fsm.StorageKey, key string, data any) error
	GetData(targetKey fsm.StorageKey, key string, to any) error
	Close() error
}

type ContextStorageWrapper struct {
	s StorageWithoutContext
}

func NewContextStorageWrapper(s StorageWithoutContext) ContextStorageWrapper {
	return ContextStorageWrapper{s: s}
}

func (c ContextStorageWrapper) GetState(_ context.Context, targetKey fsm.StorageKey) (fsm.State, error) {
	return c.s.GetState(targetKey)
}

func (c ContextStorageWrapper) SetState(_ context.Context, targetKey fsm.StorageKey, state fsm.State) error {
	return c.s.SetState(targetKey, state)
}

func (c ContextStorageWrapper) ResetState(_ context.Context, targetKey fsm.StorageKey, withData bool) error {
	return c.s.ResetState(targetKey, withData)
}

func (c ContextStorageWrapper) UpdateData(_ context.Context, targetKey fsm.StorageKey, key string, data any) error {
	return c.s.UpdateData(targetKey, key, data)
}

func (c ContextStorageWrapper) GetData(_ context.Context, targetKey fsm.StorageKey, key string, to any) error {
	return c.s.GetData(targetKey, key, to)
}

func (c ContextStorageWrapper) Close() error {
	return c.s.Close()
}
