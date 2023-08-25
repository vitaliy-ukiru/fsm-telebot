package storages

import (
	"fmt"
	"reflect"
)

type ErrWrongTypeAssign struct {
	Expect reflect.Type
	Got    reflect.Type
}

func (e ErrWrongTypeAssign) Error() string {
	return fmt.Sprintf("wrong types, can't assign %s to %s", e.Expect, e.Got)
}
