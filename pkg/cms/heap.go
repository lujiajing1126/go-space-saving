package cms

import (
	"container/heap"
)

var _ heap.Interface = (*counterList)(nil)

type counterList []*counter

type counter struct {
	count uint64
	value interface{}
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int
}

func (c *counter) Val() interface{} {
	return c.value
}

func (c counterList) Len() int {
	return len(c)
}

func (c counterList) Less(i, j int) bool {
	return c[i].count < c[j].count
}

func (c counterList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
	c[i].index = i
	c[j].index = j
}

func (c *counterList) Push(x interface{}) {
	n := len(*c)
	item := x.(*counter)
	item.index = n
	*c = append(*c, item)
}

func (c *counterList) Pop() interface{} {
	old := *c
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*c = old[0 : n-1]
	return item
}

type cmsHeap struct {
	maxLen      uint
	counterList *counterList
	cache       map[interface{}]*counter
}

func newCmsHeap(maxLen uint) *cmsHeap {
	cts := make(counterList, 0)
	heap.Init(&cts)
	return &cmsHeap{
		maxLen:      maxLen,
		counterList: &cts,
		cache:       make(map[interface{}]*counter),
	}
}

func (h *cmsHeap) Push(item interface{}, num uint64) {
	// first check if the item exists in the heap
	if c, ok := h.cache[item]; ok {
		// update counter
		c.count = num
		heap.Fix(h.counterList, c.index)
	} else {
		c := &counter{
			count: num,
			value: item,
		}
		heap.Push(h.counterList, c)
		h.cache[item] = c
		if uint(h.counterList.Len()) > h.maxLen {
			poppedItem := heap.Pop(h.counterList).(*counter).value
			delete(h.cache, poppedItem)
		}
	}
}

func (h *cmsHeap) SortedArray() []*counter {
	l := h.counterList.Len()
	counters := make([]*counter, l)
	copiedArr := make([]*counter, l)
	copy(copiedArr, *h.counterList)
	for i := 0; i < l; i++ {
		counters[l-i-1] = heap.Pop(h.counterList).(*counter)
	}
	return counters
}
