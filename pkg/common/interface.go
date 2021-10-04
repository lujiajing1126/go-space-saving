package common

type Serializable interface {
	Bytes() []byte
}

type Counter interface {
	Val() interface{}
	Count() uint64
}
