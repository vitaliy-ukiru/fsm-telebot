package file

import (
	"fmt"
	"io"
	"sync"

	"github.com/vitaliy-ukiru/fsm-telebot"
)

type WriterFunc func() (io.WriteCloser, error)

// Provider saves data to files (or streams). With custom format.
// Some providers (json, gob, base64) are already implemented in
// the `provider` sub-package
type Provider interface {
	ProviderName() string
	Save(w io.Writer, data ChatsStorage) error
	Read(r io.Reader) (ChatsStorage, error)
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// chatKey represents  pair {c: chat id, u: user id}
type chatKey struct {
	c, u int64
}

func newKey(chat, user int64) chatKey {
	return chatKey{
		c: chat,
		u: user,
	}
}

type record struct {
	state fsm.State
	data  map[string]dataCache
}

// dataCache stores data in two variants.
// Decoded in loaded
// and raw.
type dataCache struct {
	// loaded decoded content from raw via provider
	// see dataCache.get in ./internal
	loaded interface{}
	// raw content from file.
	raw []byte
}

// Storage is storage based on RAM. Drops if you stop script.
type Storage struct {
	rw       sync.RWMutex
	data     map[chatKey]record
	p        Provider
	writerFn WriterFunc
}

func NewStorage(p Provider, writerFn WriterFunc) *Storage {
	return &Storage{p: p, writerFn: writerFn, data: make(map[chatKey]record)}
}

func (s *Storage) Init(r io.Reader) error {
	if r == nil {
		return nil
	}
	dump, err := s.p.Read(r)
	if err != nil {
		return err
	}
	s.reset(dump)
	return nil
}

func (s *Storage) GetState(chatId, userId int64) (fsm.State, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.data[newKey(chatId, userId)].state, nil
}

func (s *Storage) SetState(chatId, userId int64, state fsm.State) error {
	s.do(chatId, userId, func(r *record) {
		r.state = state
	})
	return nil
}

func (s *Storage) ResetState(chatId, userId int64, withData bool) error {
	s.do(chatId, userId, func(r *record) {
		r.state = ""
		if withData {
			for key := range r.data {
				delete(r.data, key)
			}
		}
	})
	return nil
}

func (s *Storage) UpdateData(chatId, userId int64, key string, data interface{}) error {
	s.do(chatId, userId, func(r *record) {
		r.updateData(key, data)
	})
	return nil
}

func (s *Storage) GetData(chatId, userId int64, key string, to interface{}) error {
	s.rw.RLock()
	defer s.rw.RUnlock()
	d, ok := s.data[newKey(chatId, userId)].data[key]
	if !ok {
		return fsm.ErrNotFound
	}

	return d.get(to, s.p)
}

// Close saves storage data to writer from writer function.
//
// Also, the method closes writer, minimum once time.
func (s *Storage) Close() error {
	w, err := s.writerFn()
	if err != nil {
		return err
	}
	defer w.Close()

	if err := s.save(w); err != nil {
		return err
	}
	return w.Close()
}

// SaveTo saves storage data to writer.
// You can use this method to create dumps in runtime.
func (s *Storage) SaveTo(w io.Writer) error {
	return s.save(w)
}

func (s *Storage) save(w io.Writer) error {
	dump, err := s.dump()
	if err != nil {
		return err
	}

	return s.p.Save(w, dump)
}

type ProviderError struct {
	ProviderType string
	Operation    string
	Err          error
}

func (e ProviderError) Unwrap() error { return e.Err }
func (e ProviderError) Error() string {
	return fmt.Sprintf("fsm-telebot/storage/file/provider: %s: %s: %v", e.ProviderType, e.Operation, e.Err)
}
