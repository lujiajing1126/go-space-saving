package space_saving

import (
	"container/list"
)

type Counter struct {
	// value holds the real counter value.
	value int64

	// error keeps track of the error bar.
	error int64

	// data is the actual data
	data interface{}

	// bucket is the node where the parent is located.
	// We save this reference for easier detachment.
	bucket *list.Element
}

func newCounter(parent *list.Element) *Counter {
	return &Counter{
		value:  0,
		error:  0,
		data:   nil,
		bucket: parent,
	}
}

func (c *Counter) Val() interface{} {
	return c.data
}

func (c *Counter) Count() int64 {
	return c.value
}
