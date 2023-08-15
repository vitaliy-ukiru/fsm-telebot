package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func Test_endpointName(t *testing.T) {

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "constant from telebot (tele.OnText)",
			arg:  tele.OnText,
			want: "OnText",
		},
		{
			name: "inlining constant from telebot (tele.OnPhoto)",
			arg:  "\aphoto",
			want: "OnPhoto",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, endpointName(tt.arg), "endpointName(%v)", tt.arg)
		})
	}
}

func TestEndpointFormat(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "telebot event",
			arg:  tele.OnText,
			want: "OnText",
		},
		{
			name: "callback data",
			arg:  (&tele.Btn{Unique: "fsm_data"}).CallbackUnique(),
			want: "CallbackUnique(fsm_data)",
		},
		{
			name: "custom handler",
			arg:  "custom_handler",
			want: "custom_handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, EndpointFormat(tt.arg), "EndpointFormat(%v)", tt.arg)
		})
	}
}
