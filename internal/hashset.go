package internal

import "github.com/vitaliy-ukiru/fsm-telebot"

type Hashset map[fsm.State]struct{}

func NewHashset() Hashset {
	return make(Hashset)
}

func NewHashsetFromSlice(states []fsm.State) Hashset {
	h := NewHashset()
	for _, state := range states {
		h.Add(state)
	}
	return h
}

func (h Hashset) Add(s fsm.State) {
	h[s] = struct{}{}
}

func (h Hashset) Has(state fsm.State) bool {
	_, ok := h[state]
	return ok
}

func (h Hashset) Delete(state fsm.State) {
	delete(h, state)
}

func (h Hashset) Size() int {
	return len(h)
}
