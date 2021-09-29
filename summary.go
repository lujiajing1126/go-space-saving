package space_saving

import (
	"container/list"
	"errors"
	"math"
)

var (
	ErrInvalidEpsilon = errors.New("epsilon must be larger than 0.0 and less or equal than 1.0")
)

type StreamSummary struct {
	// counter is the total number of elements that have been ever seen in the data stream.
	counter int64

	// capacity is the number of counters that can be held in this StreamSummary.
	// this number is determined by the epsilon given by the user.
	capacity int64

	// epsilon is the error toleration set by the user.
	// Any element, ei, with frequency fi > \epsilon * N is guaranteed to be in the StreamSummary.
	epsilon float64

	// buckets is the doubly-linked list holding various bucket.
	buckets *list.List
}

func NewStreamSummary(epsilon float64) (*StreamSummary, error) {
	if epsilon <= 0 || epsilon > 1 {
		return nil, ErrInvalidEpsilon
	}
	capacity := math.Ceil(1.0 / epsilon)
	ss := &StreamSummary{
		capacity: int64(capacity),
		counter:  0,
		epsilon:  epsilon,
		buckets:  list.New(),
	}
	ss.init()
	return ss, nil
}

func (ss *StreamSummary) init() {
	zeroBucket := newBucket(0)
	ss.buckets.PushBack(zeroBucket)
	// init zeroBucket with counters
	var i int64
	for i = 0; i < ss.capacity; i++ {
		// add new counter
	}
}
