package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/file"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/file/provider"
	tele "gopkg.in/telebot.v3"
)

const (
	stateBegin fsm.State = "begin"
	stateRight fsm.State = "right"
	stateLeft  fsm.State = "left"
	stateEnd   fsm.State = "end"
)

const helpText = `Commands:
/data <key> - Get data from storage by key
/set_data <key> <data> - Set data by key
/begin - Start state travel
/stop - Stop state travel
/state - Sends your state
/complex <float> - Save complex structure to storage
/snap - Get complex structure from storage`

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	storagePath := flag.String(
		"storage-path",
		"",
		"Path to file with data for FSM storage",
	)

	flag.Parse()
	if *storagePath == "" {
		log.Println("setup storage-path command line argument")
		flag.PrintDefaults()
		os.Exit(1)
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 3 * time.Second},
		OnError: func(err error, c tele.Context) {
			log.Printf("[ERR] %q chat=%s", err, c.Recipient())
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	fsmStorage := file.NewStorage(
		provider.PrettyJson{
			JsonSettings: provider.JsonSettings{
				Indent: "  ",
			},
		},
		file.OpenWriter(*storagePath),
	)

	{
		storageFile, err := os.Open(*storagePath)
		if err == nil {
			defer storageFile.Close()
			if err := fsmStorage.Init(storageFile); err != nil {
				log.Fatalln(err)
			}
		}
	}

	defer func(fsmStorage *file.Storage) {
		log.Println("saving storage state in Close")
		err = fsmStorage.Close()
		if err != nil {
			log.Print("close storage error: ", err)
		}
	}(fsmStorage)

	m := fsm.NewManager(bot, nil, fsmStorage, nil)
	m.Group().Handle("/help", func(c tele.Context) error {
		return c.Send(helpText)
	})
	m.Bind("/data", fsm.DefaultState, GetDataHandler)
	m.Bind("/set_data", fsm.DefaultState, SetDataHandler)

	m.Bind("/begin", fsm.DefaultState, CommandBegin)

	m.Bind("/right", stateBegin, CommandRight)
	m.Bind(tele.OnText, stateRight, RevertMessageText)

	m.Bind("/left", stateBegin, CommandLeft)
	m.Bind(tele.OnText, stateLeft, UpperMessageText)

	m.Handle(fsm.F("/end", stateRight, stateLeft), CommandEnd)
	m.Bind("/complex", fsm.AnyState, CommandComplex)
	m.Bind("/snap", fsm.AnyState, CommandSnap)

	m.Bind("/stop", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		if err := state.Finish(false); err != nil {
			return sendErr(c, err)
		}

		return c.Send("You stop state")
	})

	m.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		userState, err := state.State()
		if err != nil {
			return sendErr(c, err)
		}
		return c.Send("Your state: " + userState.GoString())
	})

	log.Println("bot stated")
	go bot.Start()

	{
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
		<-stopChan
	}

	// This operation too longs (~4 sec.)
	// Don't wait in example.
	// But in production code I think better don't use goroutine
	go bot.Stop()
	log.Println("bot stopped")

}

func GetDataHandler(c tele.Context, state fsm.Context) error {
	key := c.Message().Payload
	if key == "" {
		return c.Send("Send key after command. Like /data my_key")
	}
	var data string
	if err := state.Get(key, &data); err != nil {
		if errors.Is(err, fsm.ErrNotFound) {
			return c.Send("Not found data in my storage :(")
		}
		return sendErr(c, err)
	}
	return c.Send("Your data: <code>"+data+"</code>", tele.ModeHTML)
}

func SetDataHandler(c tele.Context, state fsm.Context) error {
	args := strings.SplitN(c.Message().Payload, " ", 2)
	if len(args) < 2 {
		return c.Send("Ops. You must input two args. Like /set_data my_key my_data")
	}

	if err := state.Update(args[0], args[1]); err != nil {
		return sendErr(c, err)
	}

	return c.Send("I save it!")
}

func CommandBegin(c tele.Context, state fsm.Context) error {
	if err := state.Set(stateBegin); err != nil {
		return sendErr(c, err)
	}

	return c.Send("Select next step: /right or /left")
}

func CommandRight(c tele.Context, state fsm.Context) error {
	if err := state.Set(stateRight); err != nil {
		return sendErr(c, err)
	}

	return c.Send("Select next step: /end\nOr send me text and I'm reverse it.")
}

func RevertMessageText(c tele.Context, _ fsm.Context) error {
	var sb strings.Builder
	text := []rune(c.Text()) // for support unicode
	for i := len(text) - 1; i >= 0; i-- {
		sb.WriteRune(text[i])
	}
	return c.Send(sb.String())
}

func CommandLeft(c tele.Context, state fsm.Context) error {
	if err := state.Set(stateLeft); err != nil {
		return sendErr(c, err)
	}

	return c.Send("Select next step: /end\nOr send me text and I'm upper it")
}

func UpperMessageText(c tele.Context, _ fsm.Context) error {
	return c.Send(strings.ToUpper(c.Text()))

}

func CommandEnd(c tele.Context, state fsm.Context) error {
	if err := state.Set(stateEnd); err != nil {
		return sendErr(c, err)
	}

	return c.Send("Select step: /right or /left\nOr send /stop for end travel.")
}

type snap struct {
	Now   time.Time `json:"now"`
	User  int64     `json:"user"`
	Input float64   `json:"input"`
}

func CommandComplex(c tele.Context, state fsm.Context) error {
	float, err := strconv.ParseFloat(c.Message().Payload, 64)
	if err != nil {
		return c.Send(err.Error())
	}
	s := snap{
		Now:   time.Now(),
		User:  c.Sender().ID,
		Input: float,
	}
	if err := state.Update("snap", s); err != nil {
		return c.Send(err.Error())
	}
	return c.Send("saved")
}

func CommandSnap(c tele.Context, state fsm.Context) error {
	var s snap
	if err := state.Get("snap", &s); err != nil {
		return c.Send(err.Error())
	}
	return c.Send(fmt.Sprintf("%+v", s))
}

func sendErr(c tele.Context, err error) error {
	return c.Send("error: " + err.Error())
}
