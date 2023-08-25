package file

import (
	"errors"
	"fmt"
	"reflect"
)

func (r *record) updateData(key string, data interface{}) {
	if r.data == nil {
		r.data = make(map[string]dataCache)
	}
	if data == nil {
		delete(r.data, key)
	} else {
		r.data[key] = dataCache{loaded: data}
	}
}

func (r *record) resetData() {
	for key := range r.data {
		delete(r.data, key)
	}
}

// do exec `call` and save modification to storage.
// It helps not to copy the code.
func (s *Storage) do(chat, user int64, call func(*record)) {
	s.rw.Lock()
	defer s.rw.Unlock()
	key := newKey(chat, user)

	r := s.data[key]
	call(&r)
	s.data[key] = r
}

// get value from data. Priority on loaded value.
func (d *dataCache) get(to interface{}, p Provider) error {
	destValue := reflect.ValueOf(to)
	if destValue.Kind() != reflect.Pointer {
		return errors.New("dest argument must be pointer")
	}

	destElem := destValue.Elem()
	destType := destElem.Type()

	if d.loaded != nil {
		vType := reflect.TypeOf(d.loaded)
		if !vType.AssignableTo(destType) {
			return fmt.Errorf("wrong types, can't assign %s to %s", vType, destType)
		}

		destElem.Set(reflect.ValueOf(d.loaded))
		return nil
	}

	if err := p.Decode(d.raw, to); err != nil {
		return err
	}

	d.loaded = destElem.Interface()
	return nil
}
