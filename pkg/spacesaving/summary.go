package spacesaving

import (
	"container/list"
	"errors"
	"github.com/lujiajing1126/go-space-saving/pkg/common"
	"math"
)

var (
	ErrInvalidEpsilon = errors.New("epsilon must be larger than 0.0 and less or equal than 1.0")
)

type StreamSummary struct {
	// counter is the total number of elements that have been ever seen in the data stream.
	counter uint64

	// capacity is the number of counters that can be held in this StreamSummary.
	// this number is determined by the epsilon given by the user.
	capacity uint64

	// epsilon is the error toleration set by the user.
	// Any element, ei, with frequency fi > \epsilon * N is guaranteed to be in the StreamSummary.
	epsilon float64

	// buckets is the doubly-linked list holding various bucket.
	// The buckets in the list is sorted by descending order.
	buckets *list.List

	// cache maps the recorded items to the LinkedList's Node wrapping the related counter.
	cache map[interface{}]*list.Element
}

func NewStreamSummary(epsilon float64) (*StreamSummary, error) {
	if epsilon <= 0 || epsilon > 1 {
		return nil, ErrInvalidEpsilon
	}
	capacity := math.Ceil(1.0 / epsilon)
	return NewStreamSummaryWithFixedCap(uint64(capacity))
}

func NewStreamSummaryWithFixedCap(capacity uint64) (*StreamSummary, error) {
	ss := &StreamSummary{
		capacity: capacity,
		counter:  0,
		epsilon:  0,
		buckets:  list.New(),
		cache:    make(map[interface{}]*list.Element),
	}
	ss.init()
	return ss, nil
}

func (ss *StreamSummary) init() {
	zeroBucket := newBucket(0)
	// init zeroBucket with counters
	var i uint64
	elem := ss.buckets.PushBack(zeroBucket)
	for i = 0; i < ss.capacity; i++ {
		// add new Counter
		zeroBucket.children.PushBack(newCounter(elem))
	}
}

func (ss *StreamSummary) Record(item interface{}) {
	// increase the Counter since we have a new event
	ss.counter++
	if counterElem, ok := ss.cache[item]; ok {
		// if we've already recorded this item
		ss.increaseCounter(counterElem)
	} else {
		// if we haven't seen this item before
		minCounterElem := ss.buckets.Back().Value.(*bucket).children.Front()
		minCounter := minCounterElem.Value.(*counter)
		// remove old from cache
		if _, ok := ss.cache[minCounter.data]; ok {
			delete(ss.cache, minCounter.data)
		}
		// set underlying data
		minCounter.data = item
		ss.increaseCounter(minCounterElem)
		if len(ss.cache) <= int(ss.capacity) {
			minCounter.error = minCounter.value
		}
	}
}

func (ss *StreamSummary) TopK(top uint) []common.Counter {
	counters := make([]common.Counter, 0, top)
	iter := NewCounterIter(ss.buckets.Front())
	for iter.HasNext() && top > 0 {
		counters = append(counters, iter.Next())
		top = top - 1
	}
	return counters
}

func (ss *StreamSummary) increaseCounter(counterElem *list.Element) {
	c := counterElem.Value.(*counter)
	bucketParentElem := c.bucket
	bucketNextElem := bucketParentElem.Prev()
	// remove the Counter from the old bucket
	bucketParent := bucketParentElem.Value.(*bucket)
	bucketParent.children.Remove(counterElem)
	// increase the Counter value
	c.value += 1
	if bucketNextElem != nil && bucketNextElem.Value != nil && c.value == bucketNextElem.Value.(*bucket).value {
		// if the nextBucket exists and the value of the bucket equals to the current Counter value
		nextBucket := bucketNextElem.Value.(*bucket)
		c.bucket = bucketNextElem
		// push the Counter at the end of the new bucket
		counterNewElem := nextBucket.children.PushBack(c)
		// update cache with the newly created node
		ss.cache[c.data] = counterNewElem
	} else {
		// if 1. the nextBucket is nil or 2. the representative value is different
		// we have to create a new bucket
		bucketNew := newBucket(c.value)
		// push the Counter at the end of the new bucket
		counterNewElem := bucketNew.children.PushBack(c)
		bucketNewElem := ss.buckets.InsertBefore(bucketNew, bucketParentElem)
		c.bucket = bucketNewElem
		// update cache with the newly created node
		ss.cache[c.data] = counterNewElem
	}

	// delete the current bucket if no child exists
	if bucketParent.children.Len() == 0 {
		ss.buckets.Remove(bucketParentElem)
	}
}
