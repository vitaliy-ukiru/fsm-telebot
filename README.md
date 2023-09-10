# fsm-telebot

![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/vitaliy-ukiru/fsm-telebot/support%2Fv1.2.x-revert?style=flat-square&label=Go)
[![Go Reference](https://pkg.go.dev/badge/github.com/vitaliy-ukiru/fsm-telebot.svg)](https://pkg.go.dev/github.com/vitaliy-ukiru/fsm-telebot@v1.2.1)
[![Go v1.2.x](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go-support-v1.2.x.yml/badge.svg)](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go-support-v1.2.x.yml)

Finite State Machine for [telebot](https://gopkg.in/telebot.v3). 
Based on [aiogram](https://github.com/aiogram/aiogram) FSM version.

It not a full implementation FSM. It just states manager for telegram bots.

## Install:
```
go get -u github.com/vitaliy-ukiru/fsm-telebot@v1.2.1
```


## Examples
<details>
<summary>simple configuration</summary>

```go
package main

import (
	"os"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
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

	// for example using memory storage
	// but prefer will use redis or file storage.
	storage := memory.NewStorage()
	manager := fsm.NewManager(
		bot,     // tele.Bot
		nil,     // handlers will setups to this group. Default: creates new
		storage, // storage for states and data
		nil,     // context maker. Default: NewFSMContext
	)
	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		userState, err := state.State()
		if err != nil {
			return c.Send("error: " + err.Error())
		}

		return c.Send(userState.GoString())
	})

}

```

</details>

Many complex examples in directory [examples](./examples).

