package spacesaving

import (
	"container/list"
)

type CounterIterator interface {
	HasNext() bool
	Next() *Counter
}

var _ CounterIterator = (*counterIter)(nil)

type counterIter struct {
	bucketElem  *list.Element
	counterElem *list.Element
}

func NewCounterIter(bucketElem *list.Element) CounterIterator {
	return &counterIter{
		bucketElem:  bucketElem,
		counterElem: nil,
	}
}

func (c *counterIter) HasNext() bool {
	if c.bucketElem != nil && c.counterElem == nil {
		c.counterElem = c.bucketElem.Value.(*bucket).children.Front()
		return true
	}

	// we have next element in the inner chain
	if c.counterElem.Next() != nil {
		c.counterElem = c.counterElem.Next()
		return true
	} else {
		// if the next element is null,
		// reset the bucket counter if we've exhausted the current inner iterator
		if c.bucketElem.Next() != nil {
			c.bucketElem = c.bucketElem.Next()
			c.counterElem = c.bucketElem.Value.(*bucket).children.Front()
			return true
		}
	}

	return false
}

func (c *counterIter) Next() *Counter {
	return c.counterElem.Value.(*Counter)
}
