# fsm-telebot

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vitaliy-ukiru/fsm-telebot/v2.x?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/vitaliy-ukiru/fsm-telebot/v2.svg)](https://pkg.go.dev/github.com/vitaliy-ukiru/fsm-telebot/v2)
[![Go](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go.yml/badge.svg?branch=master&style=flat-square)](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/go.yml)
[![golangci-lint](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/vitaliy-ukiru/fsm-telebot/actions/workflows/golangci-lint.yml)


# Overview
Finite State Machine for [telebot](https://gopkg.in/telebot.v3). 
Based on [aiogram](https://github.com/aiogram/aiogram) FSM version.

It not a full implementation FSM. It just states manager for telegram bots.


This module build over 
[github.com/vitaliy-ukiru/telebot-filter](https://pkg.go.dev/github.com/vitaliy-ukiru/telebot-filter).




# Install
```
go get github.com/vitaliy-ukiru/fsm-telebot/v2
```

# Quick Start
## Create telebot.Bot instance
```go
bot, err := tele.NewBot(tele.Settings{
		Token:   os.Getenv("BOT_TOKEN"),
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: *debug,
})
```

## Optionally create group (recommended)
```go
g = bot.Group()
```

## Init FSM manager
```
m := fsm.New(memory.NewStorage())
```
Also, you can add optional settings. 
See [fsmopt/manager.go](./fsmopt/manager.go)

### List of settings
- fsmopt.Strategy - Setup custom strategy for user's targeting in storage
- fsmopt.ContextFactory - Setup custom FSM context factory
- fsmopt.FilterProcessor - Setup custom state filter processor

## Optionally setup middleware for context
This middleware needs only for cache fsm.Context 
object in telebot.Context.

```
g.Use(m.WrapContext)
```

## Bind handler via telebot-filter
You can use any flow of telebot-filter


### telebot-filter/dispatcher
#### Init dispatcher
```go
dp := dispatcher.NewDispatcher(g)
```
Optionally you can add router object from dispatcher.

#### Setups handlers
These methods work and with router object.
```go
dp.Dispatch(
    m.New(
		// List of options for fsm.HandlerConfig
        //  in fsmopt/handler.go
		...,
    ),
)

m.Bind(
    dp,
    // List of options for fsm.HandlerConfig
    //  in fsmopt/handler.go
    ...,
)


m.Handle(
    dp,
    tele.OnText, // endpoint
    MyState,     // state filter as fsm.StateMatcher 
    handler,     // handler as fsm.Handler
)
```
### telebot-filter/routing
With this flow available only Manager.New

And you don't need to pass endpoint option.
Because routing package don't look on that.
```
bot.Handle("/start", 
    routing.New(
        // don't specify the endpoint because it doesn't count.
        m.New(
            ...
        ),

        m.New(
            ...
        ),
    ),
)
```

# Examples
See examples in directory [examples](./examples).


# Storages
FSM need storage user state and data for user.

## Memory storage
This storage distributes with fsm-telebot/v2.
It just in-memory storage, that lose data on stop app

### Path: 
```
github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory
```

## Redis storage
This storage based on redis in distributes separated.

### Install
```
go get github.com/nacknime-official/fsm-telebot-redis-storage/v2
```

## Add your implementation to list
If you want to add you implementation in this list, 
then open pull request.
