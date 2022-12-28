package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"

	tele "gopkg.in/telebot.v3"
)

const (
	MyState fsm.State = "my_state" // Values must be unique else it breaks semantic
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	flag.Parse()

	bot, err := tele.NewBot(tele.Settings{
		Token:   os.Getenv("BOT_TOKEN"),
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: *debug,
	})
	if err != nil {
		log.Fatalln(err)
	}
	m := fsm.NewManager(bot, nil, memory.NewStorage())

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
			state.Set(MyState)
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
