package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_Insert(t *testing.T) {
	type testCase[T any] struct {
		name  string
		slice []T

		// zero value
		l         List[T]
		_typeZero T
	}

	tests := []testCase[int]{
		{
			name:  "one element",
			slice: []int{5},
		},
		{
			name:  "two elements",
			slice: []int{6, 9},
		},
		{
			name:  "more slice",
			slice: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := sliceByTypeValue(tt._typeZero, 0, len(tt.slice))
			values := tt.slice
			l := tt.l

			for _, v := range values {
				e := l.Insert(v)
				elements = append(elements, e)
			}

			assert.Equal(t, len(values), len(elements), "lengths of test data")
			assert.Equal(t, len(values), l.Len(), "length of slice and list")

			i := 0
			for e := l.Front(); e != nil; e = e.Next() {
				load := elements[i]
				assert.Equalf(t, load, e, "elements at slice[%d] and list", i)
				i++
			}
			assert.Equal(t, len(values), i, "iterator and values slice")
		})
	}
}

// sliceByTypeValue allocates slice elements of T
func sliceByTypeValue[T any](_ T, length, capacity int) []*Element[T] {
	return make([]*Element[T], length, capacity)
}

func TestList_Front(t *testing.T) {
	type testCase[T any] struct {
		name string
		l    List[T]
		want *Element[T]
	}
	e := &Element[int]{Value: 5}
	tests := []testCase[int]{
		{
			name: "nil empty list",
			// all zeroes as default
		},
		{
			name: "one element",
			l:    List[int]{head: e, len: 1},
			want: e,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.l.Front(), "Front()")
		})
	}
}

func TestWalkList(t *testing.T) {
	l := new(List[int])
	for i := 0; i < 10; i++ {
		l.Insert(i)
	}

	t.Run("front", func(t *testing.T) {
		i := 0
		for e := l.Front(); e != nil; e, i = e.Next(), i+1 {
			assert.Equalf(t, i, e.Value, "element at %d", i)
		}
	})

	t.Run("back", func(t *testing.T) {
		i := 9
		for e := l.Back(); e != nil; e, i = e.Prev(), i-1 {
			assert.Equalf(t, i, e.Value, "element at %d", i)
		}
	})
}
