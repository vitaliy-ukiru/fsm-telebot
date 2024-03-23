package container

type LinkedHashSet[K comparable] struct {
	set  map[K]*Element[K]
	list *List[K]
}

func NewLinkedHashSet[K comparable](items ...K) *LinkedHashSet[K] {
	set := new(LinkedHashSet[K])
	set.list = new(List[K])
	set.set = make(map[K]*Element[K])
	if len(items) > 0 {
		for _, item := range items {
			set.Add(item)
		}
	}
	return set
}

func (h LinkedHashSet[K]) Add(item K) {
	if h.Has(item) {
		return
	}

	node := h.list.Insert(item)
	h.set[item] = node
}

func (h LinkedHashSet[K]) Has(item K) bool {
	_, ok := h.set[item]
	return ok
}

type HashSetItem[T comparable] struct {
	el *Element[T]
}

func (item *HashSetItem[T]) Value() T {
	return item.el.Value
}

func (item *HashSetItem[T]) Next() *HashSetItem[T] {
	next := item.el.Next()
	if next == nil {
		return nil
	}

	return &HashSetItem[T]{el: next}
}

func (item *HashSetItem[T]) Prev() *HashSetItem[T] {
	prev := item.el.Prev()
	if prev == nil {
		return nil
	}

	return &HashSetItem[T]{el: prev}
}

func (h LinkedHashSet[K]) Item(item K) *HashSetItem[K] {
	node := h.set[item]
	if node == nil {
		return nil
	}
	return &HashSetItem[K]{el: node}
}

func (h LinkedHashSet[K]) Iterate(yield func(K) (next bool)) {
	for e := h.list.Front(); e != nil; e = e.Next() {
		if !yield(e.Value) {
			return
		}
	}
}
