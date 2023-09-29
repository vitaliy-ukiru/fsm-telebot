package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
)

func main() {

	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalln(err)
	}
	m := fsm.NewManager(bot, nil, memory.NewStorage(), nil)

	const MenuState fsm.State = "my_menu_state"

	// State implements StateMatcher
	// but this also matches any state,
	// because in backend it uses fsm.Is function
	m.Bind(tele.OnText, MenuState, EchoHandler)

	// Binds with match function.
	// This not have checks to any state,
	// because all logic contains in match func.
	m.BindFunc(tele.OnText, OnInputState, func(c tele.Context, state fsm.Context) error {
		return c.Send("You input: " + c.Text())
	})

	// It very simple example of state with expiration.
	// Real code will more complex, and it needs more checks
	// when just prefix and time.
	// May be it bad example.
	const ExpiredState fsm.State = "state_with_expiration"
	m.Bind("/timer", fsm.DefaultState, func(c tele.Context, state fsm.Context) error {
		payload := c.Message().Payload
		if payload == "" {
			return c.Send("Input expired time as argument in time.Duration format.\n"+
				"Format: /timer duration\n"+
				"Example <code>/timer 10s</code>", tele.ModeHTML)
		}

		expiredAfter, err := time.ParseDuration(payload)
		if err != nil {
			return c.Send("Invalid expiration time")
		}

		newState := NewTimeState(ExpiredState, time.Now().Add(expiredAfter), "")
		if err := state.Set(newState); err != nil {
			return c.Send("Fail set state: " + err.Error())
		}

		return c.Send("Send any text message after expiration time!")
	})

	m.On(tele.OnText, NewTimeMatcher(ExpiredState, ""), func(c tele.Context, state fsm.Context) error {
		return c.Send("You activated in time")
	})

	m.BindFunc(tele.OnText, PrefixMatch(ExpiredState), func(c tele.Context, state fsm.Context) error {
		defer state.Finish(true)
		return c.Send("Time is expired!")
	})
	bot.Start()

}

func OnInputState(state fsm.State) bool {
	return strings.Contains(string(state), "input")
}

func PrefixMatch(prefix fsm.State) fsm.StateMatchFunc {
	return func(state fsm.State) bool {
		return strings.Contains(string(state), string(prefix))
	}
}

const defaultSep = "|"

func NewTimeState(base fsm.State, before time.Time, sep string) fsm.State {
	if sep == "" {
		sep = defaultSep
	}

	beforeFormatted := before.Format(time.DateTime)
	return base + fsm.State(sep+beforeFormatted)
}

// TimeMatcher matches only states for which the activation time has not expired.
type TimeMatcher struct {
	prefixState string
	sep         string
}

func NewTimeMatcher(prefixState fsm.State, sep string) *TimeMatcher {
	return &TimeMatcher{prefixState: string(prefixState), sep: sep}
}

func (t *TimeMatcher) MatchState(state fsm.State) bool {
	if t.sep == "" {
		t.sep = defaultSep
	}

	parts := strings.SplitN(string(state), t.sep, 2)
	if len(parts) < 2 {
		return false
	}

	{
		// can pass empty prefix
		// or it needs equals
		if t.prefixState != "" && !strings.HasPrefix(parts[0], t.prefixState) {
			return false
		}
	}

	expireAt, err := time.Parse(time.DateTime, parts[1])
	if err != nil {
		return false
	}

	return time.Now().Before(expireAt)
}

func EchoHandler(c tele.Context, _ fsm.Context) error {
	return c.Send(c.Text())
}
