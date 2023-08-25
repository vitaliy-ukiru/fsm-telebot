package provider

import (
	"io"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
)

type Callbacks struct {
	OnSave   func(p file.Provider, w io.Writer, data file.ChatsStorage) error
	OnRead   func(p file.Provider, r io.Reader) (file.ChatsStorage, error)
	OnEncode func(p file.Provider, v interface{}) ([]byte, error)
	OnDecode func(p file.Provider, data []byte, v interface{}) error
}

// Middleware adds functional for middlewares in file.Provider
// Middlewares callbacks controls flow.
// It means what in callback you execute provider function
// Example
//
//	var p file.Provider
//	m := NewMiddleware(p, Callbacks{
//		OnEncode: func(p file.Provider, v interface{}) ([]byte, error) {
//			// actions before encoding
//			data, err := p.Encode(v)
//			// actions after encoding
//			// for example
//			log.Printf("encode value of type: %T; result=%s, err=%v", v, data, err)
//			return data, err
//		},
//	})
//
// If callback not exists for method will call raw provider.
// ADDED AS EXPERIMENT, PREFER DON'T USE IN PRODUCTION.
type Middleware struct {
	p file.Provider
	c Callbacks
}

func NewMiddleware(p file.Provider, c Callbacks) *Middleware {
	return &Middleware{p: p, c: c}
}

// Merge only non nil callbacks.
func (m *Middleware) Merge(c Callbacks) {
	if c.OnSave != nil {
		m.c.OnSave = c.OnSave
	}

	if c.OnRead != nil {
		m.c.OnRead = c.OnRead
	}

	if c.OnEncode != nil {
		m.c.OnEncode = c.OnEncode
	}

	if c.OnDecode != nil {
		m.c.OnDecode = c.OnDecode
	}
}

func (m Middleware) ProviderName() string {
	return m.p.ProviderName()
}

func (m Middleware) Save(w io.Writer, data file.ChatsStorage) error {
	fn := file.Provider.Save
	if c := m.c.OnSave; c != nil {
		fn = c
	}
	return fn(m.p, w, data)
}

func (m Middleware) Read(r io.Reader) (file.ChatsStorage, error) {
	fn := file.Provider.Read
	if c := m.c.OnRead; c != nil {
		fn = c
	}
	return fn(m.p, r)
}

func (m Middleware) Encode(v interface{}) ([]byte, error) {
	fn := file.Provider.Encode
	if c := m.c.OnEncode; c != nil {
		fn = c
	}
	return fn(m.p, v)
}

func (m Middleware) Decode(data []byte, v interface{}) error {
	fn := file.Provider.Decode
	if c := m.c.OnDecode; c != nil {
		fn = c
	}
	return fn(m.p, data, v)
}
