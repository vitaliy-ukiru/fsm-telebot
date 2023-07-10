package fsm

type hashset map[State]struct{}

func newHashset() hashset {
	return make(hashset)
}

func newHashsetFromSlice(states []State) hashset {
	h := newHashset()
	for _, state := range states {
		h.Add(state)
	}
	return h
}

func (h hashset) Add(s State) {
	h[s] = struct{}{}
}

func (h hashset) Has(state State) bool {
	_, ok := h[state]
	return ok
}

func (h hashset) Delete(state State) {
	delete(h, state)
}

func (h hashset) Size() int {
	return len(h)
}
