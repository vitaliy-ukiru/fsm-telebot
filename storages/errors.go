package storages

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotPointer = errors.New("fsm/storage: dest argument must be pointer")
var ErrInvalidValue = errors.New("fsm/storage: dest value is nil or invalid")

type ErrWrongTypeAssign struct {
	Expect reflect.Type
	Got    reflect.Type
}

func (e ErrWrongTypeAssign) Error() string {
	return fmt.Sprintf("fsm/storage: wrong types, can't assign %s to %s", e.Expect, e.Got)
}
