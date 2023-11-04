package ext

import (
	"errors"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

type Context interface {
	fsm.Context
	NextState() error
	PrevState() error
}

type Step struct {
	On       any
	State    fsm.State
	Callback Handler
}

func NewDialog(steps ...Step) *Dialog {
	return &Dialog{steps: steps}
}

func (d *Dialog) Bind(m *fsm.Manager) {
	for i, step := range d.steps {
		m.Bind(step.On, step.State, d.handler(i))
	}
}

func Route(m *fsm.Manager, steps ...Step) {
	d := Dialog{steps}
	d.Bind(m)
}

type Handler func(c tele.Context, ext Context) error

var ErrorFinished = errors.New("chain finished")

type Dialog struct {
	steps []Step
}

func (d *Dialog) handler(stepIdx int) fsm.Handler {
	step := d.steps[stepIdx]
	return func(teleCtx tele.Context, fsmCtx fsm.Context) error {
		extCtx := d.newContext(stepIdx, teleCtx, fsmCtx)

		err := step.Callback(teleCtx, extCtx)
		return err
	}
}

func (d *Dialog) newContext(stepIdx int, c tele.Context, fsmCtx fsm.Context) Context {
	return &extContext{
		Context:     fsmCtx,
		teleCtx:     c,
		dialog:      d,
		currentStep: stepIdx,
	}
}

type extContext struct {
	fsm.Context
	teleCtx     tele.Context
	dialog      *Dialog
	currentStep int
}

func (e *extContext) NextState() error {
	currentStep := e.currentStep
	if currentStep >= len(e.dialog.steps) {
		return ErrorFinished
	}

	step := e.dialog.steps[currentStep+1]
	return e.Context.Set(step.State)
}

func (e *extContext) PrevState() error {
	currentStep := e.currentStep
	if currentStep <= 0 {
		return ErrorFinished
	}

	step := e.dialog.steps[currentStep-1]
	return e.Context.Set(step.State)
}
