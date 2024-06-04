package fsm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v3"
)

func TestDefaultFilterProcessor(t *testing.T) {
	t.Run("error from State", func(t *testing.T) {
		var caughtError error

		testError := errors.New("test error")
		fsm := new(MockContext)

		fsm.EXPECT().
			State(context.Background()).
			Return(DefaultState, testError).
			Once()

		bot, err := tele.NewBot(tele.Settings{
			Offline: true,
			OnError: func(err error, context tele.Context) {
				assert.ErrorIs(t, err, testError)
				caughtError = err
			},
		})
		assert.NoError(t, err, "creating new bot")

		ctx := bot.NewContext(U)

		status := DefaultFilterProcessor(ctx, fsm, StateFilter(func(state State) bool {
			assert.Failf(
				t,
				"Unexpected call of state metcher",
				"State matcher must be not called",
			)
			return false
		}))
		assert.False(t, status, "DefaultFilterProcess")
		assert.ErrorIs(t, caughtError, testError, "Error from OnError handler")
	})
	t.Run("success process", func(t *testing.T) {
		ctx := B.NewContext(U)
		fsm := new(MockContext)
		testState := State("test_state")

		fsm.EXPECT().
			State(context.Background()).
			Return(testState, nil).
			Once()

		isMatcherCalled := false
		status := DefaultFilterProcessor(
			ctx,
			fsm,
			StateFilter(func(state State) bool {
				isMatcherCalled = true
				return AnyState.MatchState(state)
			}),
		)
		assert.True(t, status, "DefaultFilterProcessor")
		assert.True(t, isMatcherCalled, "call state matcher")
	})
}
