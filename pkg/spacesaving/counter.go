package spacesaving

import (
	"container/list"

	"github.com/lujiajing1126/go-space-saving/pkg/common"
)

var _ common.Counter = (*counter)(nil)

type counter struct {
	// value holds the real counter value.
	value uint64

	// error keeps track of the error bar.
	error uint64

	// data is the actual data
	data interface{}

	// bucket is the node where the parent is located.
	// We save this reference for easier detachment.
	bucket *list.Element
}

func newCounter(parent *list.Element) *counter {
	return &counter{
		value:  0,
		error:  0,
		data:   nil,
		bucket: parent,
	}
}

func (c *counter) Val() interface{} {
	return c.data
}

func (c *counter) Count() uint64 {
	return c.value
}
