// Package memory contains base in-memory storage.
package memory

import (
	"reflect"
	"sync"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages"
)

type chatKey struct {
	c int64 // c is Chat ID
	u int64 // u is User ID
}

// record in storage
type record struct {
	state fsm.State
	data  map[string]any
}

func newKey(chat, user int64) chatKey {
	return chatKey{
		c: chat,
		u: user,
	}
}

// Storage is storage based on RAM. Drops if you stop script.
type Storage struct {
	l       sync.RWMutex
	storage map[chatKey]record
}

// NewStorage returns new storage in memory.
func NewStorage() *Storage {
	return &Storage{
		storage: make(map[chatKey]record),
	}
}

// do exec `call` and save modification to storage.
// It helps not to copy the code.
func (m *Storage) do(chat, user int64, call func(*record)) {
	m.l.Lock()
	defer m.l.Unlock()
	key := newKey(chat, user)

	r := m.storage[key]
	call(&r)
	m.storage[key] = r
}

func (m *Storage) GetState(chatId, userId int64) (fsm.State, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	key := newKey(chatId, userId)
	return m.storage[key].state, nil
}

func (m *Storage) SetState(chatId, userId int64, state fsm.State) error {
	m.do(chatId, userId, func(r *record) {
		r.state = state
	})
	return nil
}

func (m *Storage) ResetState(chatId, userId int64, withData bool) error {
	m.do(chatId, userId, func(r *record) {
		r.state = ""
		if withData {
			clear(r.data)
		}
	})
	return nil
}

func (m *Storage) UpdateData(chatId, userId int64, key string, data any) error {
	m.do(chatId, userId, func(r *record) {
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

func (m *Storage) GetData(chatId, userId int64, key string, to any) error {
	m.l.RLock()
	defer m.l.RUnlock()
	v, ok := m.storage[newKey(chatId, userId)].data[key]
	if !ok {
		return fsm.ErrNotFound
	}

	destValue := reflect.ValueOf(to)
	if destValue.Kind() != reflect.Ptr {
		return storages.ErrNotPointer
	}
	if destValue.IsNil() || !destValue.IsValid() {
		return storages.ErrInvalidValue
	}

	destElem := destValue.Elem()
	if !destElem.IsValid() {
		return storages.ErrNotPointer
	}

	destType := destElem.Type()

	vType := reflect.TypeOf(v)
	if !vType.AssignableTo(destType) {
		return &storages.ErrWrongTypeAssign{
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
