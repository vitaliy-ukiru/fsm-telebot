package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/dispatcher"
	tf "github.com/vitaliy-ukiru/telebot-filter/telefilter"
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
	dp := dispatcher.NewDispatcher(bot.Group())
	m := fsm.New(
		memory.NewStorage(),
		fsm.StrategyDefault,
		nil,
	)
	dp.Dispatch(
		m.New(
			fsmopt.On("/stop"),            // set endpoint
			fsmopt.OnStates(fsm.AnyState), // set state filter
			fsmopt.Do(func(c tele.Context, state fsm.Context) error { // set handler
				_ = state.Finish(context.TODO(), c.Data() != "")
				return c.Send("finish")
			}),
		),
	)
	// It also for any states. Because FSM don't filter this handler
	dp.Handle("/stop", tf.RawHandler{
		Callback: m.Adapt(func(c tele.Context, state fsm.Context) error {
			s, err := state.State(context.TODO())
			if err != nil {
				return c.Send(fmt.Sprintf("can't get state: %s", err))
			}
			return c.Send("your state: " + s.GoString())
		}),
	})

	m.Bind(
		dp,
		fsmopt.OnStates(), // will handler on default state
		fsmopt.Do(func(c tele.Context, state fsm.Context) error {
			state.SetState(context.TODO(), MyState)
			_ = state.Update(context.TODO(), "payload", time.Now().Format(time.RFC850))
			return c.Send("set state")
		}),
	)
	m.Handle(
		dp,
		tele.OnText,
		MyState,
		func(c tele.Context, state fsm.Context) error {
			var payload string
			state.Data(context.TODO(), "payload", &payload)

			newPayload := time.Now().Format(time.RFC850) + "  " + c.Text()
			_ = state.Update(context.TODO(), "payload", newPayload)
			return c.Send("prev payload: " + payload)
		},
	)

	bot.Start()
}
