package file

import (
	"reflect"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/pkg/storage"
)

// do exec `call` and save modification to storage.
// It helps not to copy the code.
func (s *Storage) do(key fsm.StorageKey, call func(*record)) {
	s.rw.Lock()
	defer s.rw.Unlock()

	r := s.data[key]
	call(&r)
	s.data[key] = r
}

// get value from data. Priority on loaded value.
func (d *dataCache) get(to any, p Provider) error {
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

	if d.loaded == nil {
		if err := p.Decode(d.raw, to); err != nil {
			return err
		}
		d.loaded = destElem.Interface()
		return nil
	}

	vType := reflect.TypeOf(d.loaded)
	if !vType.AssignableTo(destType) {
		return &storage.ErrWrongTypeAssign{Expect: vType, Got: destType}
	}

	destElem.Set(reflect.ValueOf(d.loaded))
	return nil
}
