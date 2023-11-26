package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

var (
	U tele.Update
	B tele.Bot
)

func Test_wrapperContext_Get(t *testing.T) {
	teleCtx := B.NewContext(U)
	teleCtx.Set("_fsm_", 56)

	fsmCtx := Context(&fsmContext{})

	w := &wrapperContext{
		Context: teleCtx,
		fsmCtx:  fsmCtx,
	}

	tests := []struct {
		name string
		key  string
		want any
	}{
		{
			name: "fsm context key",
			key:  fsmInternalKey,
			want: fsmCtx,
		},
		{
			name: "base context key",
			key:  "_fsm_",
			want: 56,
		},
		{
			name: "unknown key",
			key:  "foo",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, w.Get(tt.key), "Get(%v)", tt.key)
		})
	}
}

func Test_tryUnwrapContext(t *testing.T) {
	type args struct {
		c tele.Context
	}
	teleCtx := B.NewContext(U)
	fsmCtx := Context(&fsmContext{})

	teleCtx.Set(fsmInternalKey, fsmCtx)

	tests := []struct {
		name  string
		args  args
		want  Context
		want1 bool
	}{
		{
			name:  "wrapped context",
			args:  args{&wrapperContext{teleCtx, fsmCtx}},
			want:  fsmCtx,
			want1: true,
		},
		{
			name:  "exist key of context",
			args:  args{teleCtx},
			want:  fsmCtx,
			want1: true,
		},
		{
			name:  "incorrect context",
			args:  args{B.NewContext(U)}, // empty context
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tryUnwrapContext(tt.args.c)
			assert.Equalf(t, tt.want, got, "tryUnwrapContext(%v)", tt.args.c)
			assert.Equalf(t, tt.want1, got1, "tryUnwrapContext(%v)", tt.args.c)
		})
	}
}
