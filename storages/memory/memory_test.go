package memory

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages"
)

func TestStorage_GetData(t *testing.T) {
	const (
		c int64 = 1 // chat id
		u int64 = 1 // user id
	)

	m := map[string]any{
		"age":   23,
		"right": true,
		"foo":   "bar",
	}

	type args struct {
		key string
		to  any
	}
	tests := []struct {
		name    string
		data    map[string]any
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "correct case",
			data:    m,
			args:    args{"age", new(int)},
			wantErr: assert.NoError,
		},
		{
			name: "not found",
			data: m,
			args: args{key: "unknown key"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, fsm.ErrNotFound, i...)
			},
		},
		{
			name: "not pointer",
			data: m,
			args: args{"right", false},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, storages.ErrNotPointer, i...)
			},
		},
		{
			name: "wrong types",
			data: m,
			args: args{"foo", new(byte)},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var e *storages.ErrWrongTypeAssign
				if !assert.ErrorAs(t, err, &e, i...) {
					return false
				}

				return assert.Equal(t, reflect.String, e.Expect.Kind(), "want string") &&
					assert.Equal(t, reflect.Uint8, e.Got.Kind(), "want byte")
			},
		},
		{
			name: "nil pointer",
			data: m,
			args: args{"age", (*int)(nil)},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, storages.ErrInvalidValue, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			m := &Storage{
				storage: map[fsm.StorageKey]record{
					{ChatID: c, UserID: u}: {
						data: tt.data,
					},
				},
			}
			tt.wantErr(
				t,
				m.GetData(
					ctx,
					fsm.StorageKey{ChatID: c, UserID: u},
					tt.args.key,
					tt.args.to,
				),
				fmt.Sprintf("GetData(%v, %v, %v, %v, %v)", ctx, c, u, tt.args.key, tt.args.to),
			)
		})
	}
}
