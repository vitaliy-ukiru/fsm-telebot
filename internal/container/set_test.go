package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSetFromSlice(t *testing.T) {
	type args[K comparable] struct {
		items []K
	}
	type testCase[K comparable] struct {
		name string
		args args[K]
		want HashSet[K]
	}
	tests := []testCase[int]{
		{
			name: "vanilla",
			args: args[int]{[]int{1, 2, 3}},
			want: HashSet[int]{
				1: {},
				2: {},
				3: {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, HashSetFromSlice(tt.args.items), "HashSetFromSlice(%v)", tt.args.items)
		})
	}
}

func TestHashSet_Has(t *testing.T) {

	type testCase[K comparable] struct {
		name string
		h    HashSet[K]
		arg  K
		want bool
	}
	tests := []testCase[int]{
		{
			name: "exists key",
			h:    HashSet[int]{5: {}},
			arg:  5,
			want: true,
		},
		{
			name: "not exists",
			h:    HashSet[int]{1: struct{}{}},
			arg:  5,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.h.Has(tt.arg), "Has(%v)", tt.arg)
		})
	}
}
