# fsm-telebot

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vitaliy-ukiru/fsm-telebot?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/vitaliy-ukiru/fsm-telebot.svg)](https://pkg.go.dev/github.com/vitaliy-ukiru/fsm-telebot)
[![Go](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go.yml/badge.svg?branch=master&style=flat-square)](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go.yml)
[![golangci-lint](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/golangci-lint.yml)

Finite State Machine for [telebot](https://gopkg.in/telebot.v3). 
Based on [aiogram](https://github.com/aiogram/aiogram) FSM version.

It not a full implementation FSM. It just states manager for telegram bots.

## Install:
```
go get github.com/vitaliy-ukiru/fsm-telebot/v2
```



## Examples
<details>
<summary>simple configuration</summary>

```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"github.com/vitaliy-ukiru/telebot-filter/dispatcher"
	tele "gopkg.in/telebot.v3"
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 3 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	dp := dispatcher.NewDispatcher(bot.Group())

	// for example using memory storage
	// but prefer will use redis or file storage.
	storage := memory.NewStorage()
	manager := fsm.NewManager(
		storage,             // storage for states and data
		fsm.StrategyDefault, // strategy handling target for storage
		nil,                 // context factory. Default: NewFSMContext
	)
	manager.Handle(dp, "/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		userState, err := state.State(context.Background())
		if err != nil {
			return c.Send("error: " + err.Error())
		}

		return c.Send(userState.GoString())
	})

}

```

</details>

Many complex examples in directory [examples](./examples).

