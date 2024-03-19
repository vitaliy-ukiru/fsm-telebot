// Package memory contains base in-memory storage.
package memory

import (
	"context"
	"reflect"
	"sync"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage"
)

var _ fsm.Storage = (*Storage)(nil)

// Storage is storage based on RAM. Drops if you stop script.
type Storage struct {
	l       sync.RWMutex
	storage map[fsm.StorageKey]record
}

// NewStorage returns new storage in memory.
func NewStorage() *Storage {
	return &Storage{
		storage: make(map[fsm.StorageKey]record),
	}
}

// record in storage
type record struct {
	state fsm.State
	data  map[string]any
}

// do exec `call` and save modification to storage.
// It helps not to copy the code.
func (m *Storage) do(key fsm.StorageKey, call func(*record)) {
	m.l.Lock()
	defer m.l.Unlock()

	r := m.storage[key]
	call(&r)
	m.storage[key] = r
}

func (m *Storage) GetState(_ context.Context, key fsm.StorageKey) (fsm.State, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.storage[key].state, nil
}

func (m *Storage) SetState(_ context.Context, key fsm.StorageKey, state fsm.State) error {
	m.do(key, func(r *record) {
		r.state = state
	})
	return nil
}

func (m *Storage) ResetState(_ context.Context, key fsm.StorageKey, withData bool) error {
	m.do(key, func(r *record) {
		r.state = ""
		if withData {
			for key := range r.data {
				delete(r.data, key)
			}
		}
	})
	return nil
}

func (m *Storage) UpdateData(_ context.Context, target fsm.StorageKey, key string, data any) error {
	m.do(target, func(r *record) {
		if r.data == nil {
			r.data = make(map[string]any)
		}
		if data == nil {
			delete(r.data, key)
		} else {
			r.data[key] = data
		}
	})
	return nil
}

func (m *Storage) GetData(_ context.Context, target fsm.StorageKey, key string, to any) error {
	m.l.RLock()
	defer m.l.RUnlock()

	r := m.storage[target]
	v, ok := r.data[key]
	if !ok {
		return fsm.ErrNotFound
	}

	destValue := reflect.ValueOf(to)
	if destValue.Kind() != reflect.Ptr {
		return storage.ErrNotPointer
	}
	if destValue.IsNil() || !destValue.IsValid() {
		return storage.ErrInvalidValue
	}

	destElem := destValue.Elem()
	if !destElem.IsValid() {
		return storage.ErrNotPointer
	}

	destType := destElem.Type()

	vType := reflect.TypeOf(v)
	if !vType.AssignableTo(destType) {
		return &storage.ErrWrongTypeAssign{
			Expect: vType,
			Got:    destType,
		}
	}
	destElem.Set(reflect.ValueOf(v))

	return nil
}

func (m *Storage) Close() error {
	return nil
}
