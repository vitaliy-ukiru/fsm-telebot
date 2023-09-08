package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategy_apply(t *testing.T) {
	const (
		c int64 = 66
		u int64 = 88
	)
	type args struct {
		chat int64
		user int64
	}
	type want struct {
		chat int64
		user int64
	}
	tests := []struct {
		name string
		s    Strategy
		args args
		want want
	}{
		{
			name: "Default",
			s:    Default,
			args: args{c, u},
			want: want{c, u},
		},
		{
			name: "OnlyUser",
			s:    OnlyUser,
			args: args{c, u},
			want: want{0, u},
		},
		{
			name: "OnlyChat",
			s:    OnlyChat,
			args: args{c, u},
			want: want{c, 0},
		},
		{
			name: "unknown strategy (as default)",
			s:    -1,
			args: args{c, u},
			want: want{c, u},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.apply(tt.args.chat, tt.args.user)
			assert.Equalf(t, tt.want.chat, got, "apply(%v, %v)", tt.args.chat, tt.args.user)
			assert.Equalf(t, tt.want.user, got1, "apply(%v, %v)", tt.args.chat, tt.args.user)
		})
	}
}
