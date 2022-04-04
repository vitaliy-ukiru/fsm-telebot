package main

import (
	"log"
	"os"
	"time"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages"

	tele "gopkg.in/telebot.v3"
)

const (
	MyState fsm.State = "my_state" // Values must be unique else it breaks semantic
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:   os.Getenv("BOT_TOKEN"),
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	m := fsm.NewManager(bot.Group(), storages.NewMemoryStorage())

	// For any state
	m.Bind("/stop", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		_ = state.Finish(c.Data() != "")
		return c.Send("finish")
	})
	// It also for any states. Because manager don't filter this handler
	bot.Handle("/state",
		m.HandlerAdapter(func(c tele.Context, state fsm.FSMContext) error {
			return c.Send("your state: " + state.State().String())
		}),
	)

	bot.Handle("/set", m.ForState(fsm.DefaultState,
		func(c tele.Context, state fsm.FSMContext) error {
			_ = state.Set(MyState)
			_ = state.Update("payload", time.Now().Format(time.RFC850))
			return c.Send("set state")
		},
	))

	m.Handle(fsm.F(tele.OnText, MyState),
		func(c tele.Context, state fsm.FSMContext) error {
			payload, _ := state.Get("payload")
			_ = state.Update("payload", time.Now().Format(time.RFC850)+"  "+c.Text())
			return c.Send("prev payload: " + (payload).(string))
		},
	)

	bot.Start()
}
