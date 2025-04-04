package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/internal/container"
)

func TestStateGroup_MatchState(t *testing.T) {

	tests := []struct {
		name string
		sg   StateGroup
		arg  State
		want bool
	}{
		{
			name: "success",
			sg:   NewStateGroup("test", "state1", "state2"),
			arg:  "test:state1",
			want: true,
		},
		{
			name: "not contains",
			sg:   NewStateGroup("test", "state1"),
			arg:  "state1", // state must don't match, because it not has prefix
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.sg.MatchState(tt.arg), "MatchState(%v)", tt.arg)
		})
	}
}

func TestStateGroup_New(t *testing.T) {
	type fields struct {
		prefix string
		group  *container.LinkedHashSet[State]
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   State
	}{
		{
			name: "add prefix",
			fields: fields{
				prefix: "test",
				group:  container.NewLinkedHashSet[State](),
			},
			args: args{
				name: "add",
			},
			want: "test:add",
		},
		{
			name: "don't add prefix",
			fields: fields{
				prefix: "test",
				group:  container.NewLinkedHashSet[State](),
			},
			args: args{
				name: "test:add",
			},
			want: "test:add",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sg := StateGroup{
				prefix: tt.fields.prefix,
				group:  tt.fields.group,
			}
			assert.Equalf(t, tt.want, sg.New(tt.args.name), "New(%v)", tt.args.name)
		})
	}
}

func TestStateGroup_Next(t *testing.T) {
	tests := []struct {
		name string
		sg   StateGroup
		arg  State
		want State
	}{
		{
			name: "empty group",
			sg:   NewStateGroup("test"),
			arg:  "test:my_state",
			want: DefaultState,
		},
		{
			name: "normal behavior",
			sg:   NewStateGroup("test", "a", "b", "c"),
			arg:  "test:b",
			want: "test:c",
		},
		{
			name: "unknown state",
			sg:   NewStateGroup("test", "a", "b"),
			arg:  "test:c",
			want: DefaultState,
		},
		{
			name: "last state",
			sg:   NewStateGroup("test", "a", "b"),
			arg:  "test:b",
			want: DefaultState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.sg.Next(tt.arg), "Next(%v)", tt.arg)
		})
	}
}

func TestStateGroup_Prev(t *testing.T) {

	tests := []struct {
		name string
		sg   StateGroup
		arg  State
		want State
	}{
		{
			name: "empty group",
			sg:   NewStateGroup("test"),
			arg:  "test:my_state",
			want: DefaultState,
		},
		{
			name: "normal behavior",
			sg:   NewStateGroup("test", "a", "b", "c"),
			arg:  "test:b",
			want: "test:a",
		},
		{
			name: "unknown state",
			sg:   NewStateGroup("test", "a", "b"),
			arg:  "test:c",
			want: DefaultState,
		},
		{
			name: "first state",
			sg:   NewStateGroup("test", "a", "b"),
			arg:  "test:a",
			want: DefaultState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.sg.Prev(tt.arg), "Prev(%v)", tt.arg)
		})
	}
}

func TestState_MatchState(t *testing.T) {
	tests := []struct {
		name  string
		state State
		other State
		want  bool
	}{
		{
			name:  "same state",
			state: "a",
			other: "a",
			want:  true,
		},
		{
			name:  "different state",
			state: "a",
			other: "b",
			want:  false,
		},
		{
			name:  "state is Any",
			state: AnyState,
			other: "test",
			want:  true,
		},
		{
			name:  "other is Any",
			state: "test",
			other: AnyState,
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.state.MatchState(tt.other), "MatchState(%v)", tt.other)
		})
	}
}

func TestState_GoString(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want string
	}{
		{
			name: "default",
			s:    DefaultState,
			want: "DefaultState",
		},
		{
			name: "any",
			s:    AnyState,
			want: "AnyState",
		},
		{
			name: "custom",
			s:    "CustomState",
			want: "State(CustomState)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.GoString(), "GoString()")
		})
	}
}
