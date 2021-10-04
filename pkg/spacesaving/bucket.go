package spacesaving

import (
	"container/list"
)

// bucket holds the Counter with the same frequency as a linked list.
// Multiple buckets are kept in a doubly linked list, sorted by their values.
type bucket struct {
	// The value of the parent bucket is the same as the counters’ value of all of its elements.
	value uint64
	// children is the LinkedList of the counters in this bucket.
	children *list.List
}

func newBucket(val uint64) *bucket {
	return &bucket{
		value:    val,
		children: list.New(),
	}
}
