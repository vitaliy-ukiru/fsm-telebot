package storage

import (
	"sync"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
)

type chatKey struct {
	// c is Chat ID
	c int64
	// u is User ID
	u int64
}

type record struct {
	state fsm.State
	data  map[string]interface{}
}

type MemoryStorage struct {
	lock    sync.RWMutex
	storage map[chatKey]record
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[chatKey]record),
	}
}

func (r *record) updateData(key string, data interface{}) {
	if r.data == nil {
		r.resetData()
	}
	if data == nil {
		delete(r.data, key)
	} else {
		r.data[key] = data
	}
}

func (r *record) resetData() {
	r.data = make(map[string]interface{})
}

func (m *MemoryStorage) do(key chatKey, call func(*record)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	r := m.storage[key]
	call(&r)
	m.storage[key] = r
}

func (m *MemoryStorage) GetState(chatId, userId int64) fsm.State {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.storage[newKey(chatId, userId)].state
}

func (m *MemoryStorage) SetState(chatId, userId int64, state fsm.State) error {
	m.do(newKey(chatId, userId), func(r *record) {
		r.state = state
	})
	return nil
}

func (m *MemoryStorage) ResetState(chatId, userId int64, withData bool) error {
	m.do(newKey(chatId, userId), func(r *record) {
		r.state = ""
		if withData {
			r.resetData()
		}
	})
	return nil
}

func (m *MemoryStorage) UpdateData(chatId, userId int64, key string, data interface{}) error {
	m.do(newKey(chatId, userId), func(r *record) {
		r.updateData(key, data)
	})
	return nil
}

func (m *MemoryStorage) GetData(chatId, userId int64, key string) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, ok := m.storage[newKey(chatId, userId)].data[key]
	if !ok {
		return nil, fsm.ErrNotFound
	}
	return v, nil
}

func (m *MemoryStorage) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storage = nil
	return nil
}

func newKey(chat, user int64) chatKey {
	return chatKey{
		c: chat,
		u: user,
	}
}
