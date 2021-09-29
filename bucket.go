package space_saving

// bucket holds the counter with the same frequency as a linked list.
// Multiple buckets are kept in a doubly linked list, sorted by their values.
type bucket struct {
	// The value of the parent bucket is the same as the counters’ value of all of its elements.
	value int64
	// ptr is the pointer to one of the counters in this bucket.
	ptr *counter
}

func newBucket(val int64) *bucket {
	return &bucket{
		value: val,
	}
}
