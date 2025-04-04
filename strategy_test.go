package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v4"
)

func Test_extractKeyWithStrategy(t *testing.T) {
	type args struct {
		c        tele.Context
		strategy Strategy
	}
	bot := tele.Bot{Me: &tele.User{ID: 10}}

	goodCtx := bot.NewContext(tele.Update{
		Message: &tele.Message{
			Chat:         &tele.Chat{ID: 1},
			Sender:       &tele.User{ID: 2},
			TopicMessage: true,
			ThreadID:     3,
		},
	})
	tests := []struct {
		name   string
		args   args
		want   StorageKey
		wantOk bool
	}{

		{
			name: "default",
			args: args{
				c:        goodCtx,
				strategy: StrategyDefault,
			},
			want:   StorageKey{BotID: 10, ChatID: 1, UserID: 2},
			wantOk: true,
		},
		{
			name: "chat only",
			args: args{
				c:        goodCtx,
				strategy: StrategyChat,
			},
			want:   StorageKey{BotID: 10, ChatID: 1, UserID: 1},
			wantOk: true,
		},

		{
			name: "global user",
			args: args{
				c:        goodCtx,
				strategy: StrategyGlobalUser,
			},
			want:   StorageKey{BotID: 10, ChatID: 2, UserID: 2},
			wantOk: true,
		},

		{
			name: "default",
			args: args{
				c:        goodCtx,
				strategy: StrategyUserInTopic,
			},
			want:   StorageKey{BotID: 10, ChatID: 1, UserID: 2, ThreadID: 3},
			wantOk: true,
		},

		{
			name: "default",
			args: args{
				c:        goodCtx,
				strategy: StrategyChatTopic,
			},
			want:   StorageKey{BotID: 10, ChatID: 1, UserID: 1, ThreadID: 3},
			wantOk: true,
		},
		{
			name: "inline query",
			args: args{
				c: bot.NewContext(tele.Update{
					Query: &tele.Query{
						Sender: &tele.User{ID: 2},
					},
				}),
				strategy: StrategyDefault,
			},
			want:   StorageKey{BotID: 10, ChatID: 2, UserID: 2},
			wantOk: true,
		},
		{
			name: "update without sender",
			args: args{
				c: bot.NewContext(tele.Update{
					Poll: &tele.Poll{},
				}),
				strategy: StrategyDefault,
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractKeyWithStrategy(tt.args.c, tt.args.strategy)
			assert.Equalf(t, tt.want, got, "extractKeyWithStrategy(%v, %s)", tt.args.c, tt.args.strategy)
			assert.Equalf(t, tt.wantOk, got1, "extractKeyWithStrategy(%v, %s)", tt.args.c, tt.args.strategy)
		})
	}
}

func TestStrategy_String(t *testing.T) {
	tests := []struct {
		name string
		s    Strategy
		want string
	}{
		{
			name: "user in topic",
			s:    StrategyUserInTopic,
			want: "StrategyUserInTopic",
		},
		{
			name: "invalid",
			s:    -1,
			want: "Strategy(-1)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.String(), "String()")
		})
	}
}
