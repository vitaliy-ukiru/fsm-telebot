package internal

import (
	"reflect"

	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage"
)

func SetupValue(value any, to any) error {
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

	vType := reflect.TypeOf(value)
	if !vType.AssignableTo(destType) {
		return &storage.ErrWrongTypeAssign{
			Expect: vType,
			Got:    destType,
		}
	}
	destElem.Set(reflect.ValueOf(value))
	return nil
}
