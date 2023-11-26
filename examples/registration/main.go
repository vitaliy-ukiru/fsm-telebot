package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const SuperuserId tele.ChatID = 986576715 // Set your tg id

var (
	InputSG           = fsm.NewStateGroup("reg")
	InputNameState    = InputSG.New("name")
	InputAgeState     = InputSG.New("age")
	InputHobbyState   = InputSG.New("hobby")
	InputConfirmState = InputSG.New("confirm")
)

var debug = flag.Bool("debug", false, "log debug info")

var (
	regBtn    = tele.Btn{Text: "📝 Start input form"}
	cancelBtn = tele.Btn{Text: "❌ Cancel Form"}

	confirmBtn      = tele.Btn{Text: "✅ Confirm and send", Unique: "confirm"}
	resetFormBtn    = tele.Btn{Text: "🔄 Reset form", Unique: "reset"}
	cancelInlineBtn = tele.Btn{Text: "❌ Cancel Form", Unique: "cancel"}
)

func main() {
	flag.Parse()

	bot, err := tele.NewBot(tele.Settings{
		Token:     os.Getenv("BOT_TOKEN"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
		Verbose:   *debug,
		OnError: func(err error, c tele.Context) {
			log.Printf("[ERR] %q chat=%s", err, c.Recipient())
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	storage := memory.NewStorage()
	defer storage.Close()

	manager := fsm.NewManager(bot, nil, storage, fsm.StrategyDefault, nil)

	bot.Use(middleware.AutoRespond())

	// commands
	bot.Handle("/start", OnStart(regBtn))
	manager.Bind("/reg", fsm.DefaultState, OnStartRegister)
	manager.Bind("/cancel", fsm.AnyState, OnCancelForm)

	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		ctx := context.TODO()
		s, err := state.State(ctx)
		if err != nil {
			return c.Send(fmt.Sprintf("can't get state: %s", err))
		}
		return c.Send(s.GoString())
	})

	// buttons
	manager.Bind(&regBtn, fsm.DefaultState, OnStartRegister)
	manager.Bind(&cancelBtn, fsm.AnyState, OnCancelForm)

	// form
	manager.Bind(tele.OnText, InputNameState, OnInputName)
	manager.Bind(tele.OnText, InputAgeState, OnInputAge)
	manager.Bind(tele.OnText, InputHobbyState, OnInputHobby)
	manager.Bind(&confirmBtn, InputConfirmState, OnInputConfirm, EditFormMessage("Now check y", "Y"))
	manager.Bind(&resetFormBtn, InputConfirmState, OnInputResetForm, EditFormMessage("Now check your", "Your old"))
	manager.Bind(&cancelInlineBtn, InputConfirmState, OnCancelForm, DeleteAfterHandler)

	log.Println("Handlers configured")
	bot.Start()
}

func OnStart(startReg tele.Btn) tele.HandlerFunc {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(startReg))
	menu.ResizeKeyboard = true

	return func(c tele.Context) error {
		log.Println("new user", c.Sender().ID)
		return c.Send(
			"<b>Welcome!</b>\n"+
				"Send /reg for start input form\n"+
				"Send /cancel for cancel input form", menu)
	}
}

func OnStartRegister(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(cancelBtn))
	menu.ResizeKeyboard = true

	ctx := context.TODO()
	state.SetState(ctx, InputNameState)
	return c.Send("Great! How your name?", menu)
}

func OnInputName(c tele.Context, state fsm.Context) error {
	name := c.Message().Text
	ctx := context.TODO()
	go state.Update(ctx, "name", name)
	go state.SetState(ctx, InputAgeState)
	return c.Send(fmt.Sprintf("Okay, %s. How old are you?", name))
}

func OnInputAge(c tele.Context, state fsm.Context) error {
	age, err := strconv.Atoi(c.Message().Text)
	if err != nil || age <= 0 || age > 200 {
		return c.Send("Incorrect age. Retry again.")
	}

	ctx := context.TODO()
	go state.Update(ctx, "age", age)
	go state.SetState(ctx, InputHobbyState)

	return c.Send("Great! What is your hobby?")
}

func OnInputHobby(c tele.Context, state fsm.Context) error {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(confirmBtn),
		m.Row(resetFormBtn, cancelInlineBtn),
	)

	ctx := context.TODO()
	go state.Update(ctx, "hobby", c.Message().Text)
	go state.SetState(ctx, InputConfirmState)

	var (
		senderName string
		senderAge  int
	)

	state.MustGet(ctx, "name", &senderName)
	state.MustGet(ctx, "age", &senderAge)

	c.Send("Wow, interesting!")
	return c.Send(fmt.Sprintf(
		"Now check your form:\n"+
			"<i>Name</i>: %q\n"+
			"<i>Age</i>: %d\n"+
			"<i>Hobby</i>: %q\n",
		senderName,
		senderAge,
		c.Message().Text,
	), m)
}

func OnInputConfirm(c tele.Context, state fsm.Context) error {
	ctx := context.TODO()
	defer state.Finish(ctx, true)
	var (
		senderName  string
		senderAge   int
		senderHobby string
	)
	state.MustGet(ctx, "name", &senderName)
	state.MustGet(ctx, "age", &senderAge)
	state.MustGet(ctx, "hobby", &senderHobby)

	data, _ := json.Marshal(tele.M{
		"name":  senderName,
		"age":   senderAge,
		"hobby": senderHobby,
	})
	log.Printf("new form: %s", data)

	var username string
	if c.Sender().Username != "" {
		username = "@" + c.Sender().Username + " " // whitespace for formatting
	}

	_, err := c.Bot().Send(SuperuserId, fmt.Sprintf(
		"New form:\n"+
			"<i>Name</i>: %q\n"+
			"<i>Age</i>: %d\n"+
			"<i>Hobby</i>: %q\n"+
			"<a href=\"tg://user?id=%d\">Sender</a> %s<code>[%d]</code>\n",
		senderName,
		senderAge,
		senderHobby,
		c.Sender().ID,
		username,
		c.Sender().ID, // sometimes links don't work due to the privacy settings

	))
	if err != nil {
		c.Bot().OnError(err, c)
	}
	return c.Send("Form accepted", tele.RemoveKeyboard)
}

func OnCancelForm(c tele.Context, state fsm.Context) error {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(regBtn))
	menu.ResizeKeyboard = true

	ctx := context.TODO()
	go state.Finish(ctx, true)
	return c.Send("Form entry canceled. Your input data has been deleted.", menu)
}

func OnInputResetForm(c tele.Context, state fsm.Context) error {
	ctx := context.TODO()
	go state.SetState(ctx, InputNameState)
	c.Send("Okay! Start again.")
	return c.Send("How your name?")
}

func EditFormMessage(old, new string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			strOffset := utf8.RuneCountInString(old)
			if nLen := utf8.RuneCountInString(new); nLen > 1 {
				strOffset -= nLen - 1
			}

			entities := make(tele.Entities, len(c.Message().Entities))
			for i, entity := range c.Message().Entities {
				entity.Offset -= strOffset
				entities[i] = entity
			}
			defer func() {
				err := c.EditOrSend(strings.Replace(c.Message().Text, old, new, 1), entities)
				if err != nil {
					c.Bot().OnError(err, c)
				}
			}()
			return next(c)
		}
	}
}

func DeleteAfterHandler(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func(c tele.Context) {
			if err := c.Delete(); err != nil {
				c.Bot().OnError(err, c)
			}
		}(c)
		return next(c)
	}
}
