package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedHashSet_Has(t *testing.T) {
	type testCase[K comparable] struct {
		name  string
		setup func() *LinkedHashSet[K]
		args  K
		want  bool
	}
	tests := []testCase[int]{
		{
			name: "exists",
			setup: func() *LinkedHashSet[int] {
				return NewLinkedHashSet(1, 23)
			},
			args: 23,
			want: true,
		},
		{
			name: "not found",
			setup: func() *LinkedHashSet[int] {
				return NewLinkedHashSet(5, 6)
			},
			args: 4,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup()
			assert.Equalf(t, tt.want, h.Has(tt.args), "Has(%v)", tt.args)
		})
	}
}

func TestLinkedHashSet_Item(t *testing.T) {

	type result[K comparable] struct {
		value K
		next  *K
		prev  *K
	}

	type testCase[K comparable] struct {
		name  string
		setup func() *LinkedHashSet[K]
		args  K
		want  *result[K]
	}

	baseSetup := func() *LinkedHashSet[int] {
		return NewLinkedHashSet(1, 2, 3)
	}

	ptr := func(x int) *int {
		return &x
	}

	tests := []testCase[int]{
		{
			name: "empty",
			setup: func() *LinkedHashSet[int] {
				return new(LinkedHashSet[int])
			},
			args: 1,
			want: nil,
		},
		{
			name: "one element",
			setup: func() *LinkedHashSet[int] {
				return NewLinkedHashSet(1)
			},
			args: 1,
			want: &result[int]{value: 1},
		},
		{
			name:  "first",
			setup: baseSetup,
			args:  1,
			want: &result[int]{
				value: 1,
				next:  ptr(2),
				prev:  nil,
			},
		},
		{
			name:  "middle",
			setup: baseSetup,
			args:  2,
			want: &result[int]{
				value: 2,
				next:  ptr(3),
				prev:  ptr(1),
			},
		},

		{
			name:  "last",
			setup: baseSetup,
			args:  3,
			want: &result[int]{
				value: 3,
				next:  nil,
				prev:  ptr(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup()

			item := h.Item(tt.args)
			if tt.want == nil {
				assert.Nilf(t, item, "Item(%v)", tt.args)
				return
			}

			assert.Equalf(t, tt.want.value, item.Value(), "Item(%v).Value()", tt.args)

			if tt.want.next == nil {
				assert.Nilf(t, item.Next(), "Item(%v).Next()", tt.args)
			} else {
				assert.Equalf(t, *tt.want.next, item.Next().Value(), "Item(%v).Next().Value()", tt.args)
			}

			if tt.want.prev == nil {
				assert.Nilf(t, item.Prev(), "Item(%v).Prev()", tt.args)
			} else {
				assert.Equalf(t, *tt.want.prev, item.Prev().Value(), "Item(%v).Prev().Value()", tt.args)
			}

			//assert.Equalf(t, tt.want, h.Item(tt.args), "Item(%v)", tt.args)
		})
	}
}

func TestLinkedHashSet_Iterate(t *testing.T) {
	type testCase[K comparable] struct {
		name  string
		setup func() *LinkedHashSet[K]
		yield func(K) (next bool)
		want  []int
	}
	tests := []testCase[int]{
		{
			name: "all items",
			setup: func() *LinkedHashSet[int] {
				return NewLinkedHashSet(1, 4, 2, 3, 4, 5)
			},
			yield: func(_ int) (next bool) { return true },
			want:  []int{1, 4, 2, 3, 5},
		},
		{
			name: "break loop",
			setup: func() *LinkedHashSet[int] {
				return NewLinkedHashSet(1, 2, 3, 4, 5)
			},
			yield: func(i int) (next bool) { return i <= 3 },
			want:  []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup()
			var buff []int

			h.Iterate(func(i int) (next bool) {
				buff = append(buff, i)
				return tt.yield(i)
			})

			assert.Equalf(t, tt.want, buff, "Iterate")
		})
	}
}
