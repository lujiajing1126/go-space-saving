package cms

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmsHeap_Push(t *testing.T) {
	tester := assert.New(t)
	h := newCmsHeap(10)
	tester.NotNil(h)
	h.Push("a", 10)
	h.Push("b", 100)
	h.Push("c", 40)
	h.Push("d", 60)
	tester.Len(h.cache, 4)
	tester.Equal((*h.counterList)[0].count, uint64(10))
}

func TestCmsHeap_Keep_Only_Top3(t *testing.T) {
	tester := assert.New(t)
	h := newCmsHeap(3)
	tester.NotNil(h)
	h.Push("a", 10)
	h.Push("b", 100)
	h.Push("c", 40)
	h.Push("d", 60)
	tester.Len(h.cache, 3)
	tester.Equal((*h.counterList)[0].count, uint64(40))
}

func TestCmsHeap_SortedArray(t *testing.T) {
	tester := assert.New(t)
	h := newCmsHeap(3)
	tester.NotNil(h)
	h.Push("a", 10)
	h.Push("b", 100)
	h.Push("c", 40)
	h.Push("d", 60)
	tester.Len(h.cache, 3)
	sortedArr := h.SortedArray()
	tester.True(sort.SliceIsSorted(sortedArr, func(i, j int) bool {
		return sortedArr[i].count > sortedArr[j].count
	}))
}
