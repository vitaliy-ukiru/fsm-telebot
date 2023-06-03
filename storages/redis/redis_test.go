package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_generateKey(t *testing.T) {
	s := &Storage{pref: StorageSettings{Prefix: "test"}}
	type args struct {
		chat    int64
		user    int64
		keyType keyType
		keys    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple key",
			args: args{
				keyType: dataKey,
				keys:    []string{"myKey"},
			},
			want: "test:0:0:data:myKey",
		},
		{
			name: "multiple key",
			args: args{
				keyType: dataKey,
				keys:    []string{"multiple", "keys"},
			},
			want: "test:0:0:data:multiple:keys",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				s.generateKey(tt.args.chat, tt.args.user, tt.args.keyType, tt.args.keys...),
				"generateKey(%v, %v, %v, %v)",
				tt.args.chat,
				tt.args.user,
				tt.args.keyType,
				tt.args.keys,
			)
		})
	}
	t.Run("key type", func(t *testing.T) {
		assert.NotEqual(
			t,
			s.generateKey(0, 0, stateKey),
			s.generateKey(0, 0, dataKey, "state"),
			"state key and data[state] equals",
		)
	})
}
