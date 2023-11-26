package fsm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

func TestManagerOneEndpoint(t *testing.T) {
	// tests have on base
	// but this testing some handlers
	// at one endpoint

	bot, _ := tele.NewBot(tele.Settings{
		OnError: func(err error, _ tele.Context) {
			assert.NoError(t, err)
		},
		Synchronous: true,
		Offline:     true,
	})

	ctxMock := &MockContext{}

	m := &Manager{
		group:        bot.Group(),
		contextMaker: func(_ Storage, _ StorageKey) Context { return ctxMock },
		handlers:     handlerMapping{},
	}

	const (
		testState1 State = "test_state_1"
		testState2 State = "test_state_2"
	)

	h := &handlerMock{ctx: map[string]tele.Context{}}

	type args struct {
		states []State
		h      Handler
		ms     []tele.MiddlewareFunc
	}

	var testCases = []struct {
		name string

		group []tele.MiddlewareFunc
		args  args

		mockState State

		wantCallFunc string
		wantCtxVars  map[string]any
	}{
		{
			name:      "first_state",
			mockState: testState1,
			args: args{
				states: slice(testState1),
				h:      h.H1,
				ms: slice(
					middlewareSetCtx("a", 1),
					middlewareSetCtx("b", 2),
				),
			},
			group: slice(middlewareSetCtx("g1", -1)),

			wantCallFunc: "H1",
			wantCtxVars: map[string]any{
				"a":  1,
				"b":  2,
				"c":  nil,
				"d":  nil,
				"g1": -1,
				"g2": nil,
			},
		},
		{
			name:      "second_state",
			mockState: testState2,
			group:     slice(middlewareSetCtx("g2", -2)),
			args: args{
				states: slice(testState2),
				h:      h.H2,
				ms: slice(
					middlewareSetCtx("c", 3),
					middlewareSetCtx("d", 4),
				),
			},

			wantCallFunc: "H2",
			wantCtxVars: map[string]any{
				"a":  nil,
				"b":  nil,
				"c":  3,
				"d":  4,
				"g1": nil,
				"g2": -2,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			call := ctxMock.EXPECT().State(context.TODO())
			call.Return(tt.mockState, nil)
			defer call.Unset()

			m.Use(tt.group...)
			defer func() {
				m.list = m.list[0:0]
			}()

			// execute handler
			{
				m.handle("test", tt.args.states, tt.args.h, tt.args.ms)
				h.On(tt.wantCallFunc).Return()

				bot.ProcessUpdate(tele.Update{
					Message: &tele.Message{Text: "test"}, // for jump to tele.OnText
				})

				h.AssertCalled(t, tt.wantCallFunc)
			}
			c := h.ctx[tt.wantCallFunc]
			require.NotNil(t, c)
			for key, want := range tt.wantCtxVars {
				assert.Equalf(t, want, c.Get(key), "tele.Context.Get(%s)", key)
			}
		})

	}
}

func middlewareSetCtx(key string, value any) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(key, value)
			return next(c)
		}
	}
}

type handlerMock struct {
	mock.Mock
	ctx map[string]tele.Context
}

func (h *handlerMock) H1(c tele.Context, _ Context) error {
	h.Called()
	h.ctx["H1"] = c
	return nil
}

func (h *handlerMock) H2(c tele.Context, _ Context) error {
	h.Called()
	h.ctx["H2"] = c
	return nil
}

func slice[T any](t ...T) []T { return t }
