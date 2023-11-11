package container

type Set[K comparable] map[K]struct{}

func HashSetFromSlice[K comparable](items []K) Set[K] {
	h := make(Set[K])
	for _, item := range items {
		h.Add(item)
	}
	return h
}

func (h Set[K]) Add(item K) {
	h[item] = struct{}{}
}

func (h Set[K]) Has(item K) bool {
	_, ok := h[item]
	return ok
}
