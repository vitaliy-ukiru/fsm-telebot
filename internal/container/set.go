package container

type HashSet[K comparable] map[K]struct{}

func HashSetFromSlice[K comparable](items []K) HashSet[K] {
	h := make(HashSet[K])
	for _, item := range items {
		h.Add(item)
	}
	return h
}

func (h HashSet[K]) Add(item K) {
	h[item] = struct{}{}
}

func (h HashSet[K]) Has(item K) bool {
	_, ok := h[item]
	return ok
}
