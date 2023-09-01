package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitaliy-ukiru/fsm-telebot/internal"
)

func Test_handlerStorage_findHandler(t *testing.T) {
	type args struct {
		endpoint string
		state    State
	}

	set := func(s ...State) internal.HashSet[State] {
		return internal.HashSetFromSlice(s)
	}
	tests := []struct {
		name     string
		handlers map[string][]handlerEntry
		args     args
		want     handlerEntry
		wantOk   bool
	}{
		{
			name: "default",
			handlers: map[string][]handlerEntry{
				"test": {
					{
						states: set("test_state"),
					},
				},
			},
			args: args{"test", "test_state"},
			want: handlerEntry{
				states: set("test_state"),
			},
			wantOk: true,
		},
		{
			name: "many handlers",
			handlers: map[string][]handlerEntry{
				"test": {

					{states: set("test_many_1")},
					{states: set("test_many_2")},
					{states: set("test_many_3")},
				},
			},
			args:   args{"test", "test_many_2"},
			want:   handlerEntry{states: set("test_many_2")},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		m := make(handlerStorage)
		for e, entries := range tt.handlers {
			for _, entry := range entries {
				m.insert(e, entry)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, got1 := m.findHandler(tt.args.endpoint, tt.args.state)
			assert.Equalf(t, tt.want, got, "findHandler(%v, %v)", tt.args.endpoint, tt.args.state)
			assert.Equalf(t, tt.wantOk, got1, "findHandler(%v, %v)", tt.args.endpoint, tt.args.state)
		})
	}
}
