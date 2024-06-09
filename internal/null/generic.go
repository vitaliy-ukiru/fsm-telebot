package null

type Nullable[T any] struct {
	Value T
	Valid bool
}

func (null *Nullable[T]) Set(value T) {
	null.Value = value
	null.Valid = true
}
