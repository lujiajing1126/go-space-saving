package cms

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"math"
)

type CountMinSketch struct {
	// w is the width of the matrix.
	w uint
	// d is the depth of the matrix.
	d uint
	// m is a two-dimension array for storing count number, aka matrix.
	m [][]uint64
	// h is the hash function.
	h hash.Hash64
	// heap
	topHeap *cmsHeap
}

func New(d uint, w uint, topK uint) (*CountMinSketch, error) {
	if d <= 0 || w <= 0 {
		return nil, errors.New("width and/or depth must be greater than zero")
	}

	s := &CountMinSketch{
		d:       d,
		w:       w,
		h:       fnv.New64(),
		m:       make([][]uint64, d),
		topHeap: newCmsHeap(topK),
	}
	for i := range s.m {
		s.m[i] = make([]uint64, s.w)
	}
	return s, nil
}

// NewWithEstimates creates a new Count-Min Sketch with given error rate and confidence.
// Accuracy guarantees will be made in terms of a pair of user specified parameters,
// ε and δ, meaning that the error in answering a query is within a factor of ε with
// probability δ
func NewWithEstimates(epsilon, delta float64, topK uint) (*CountMinSketch, error) {
	if epsilon <= 0 || epsilon >= 1 {
		return nil, errors.New("value of epsilon should be in range of (0, 1)")
	}
	if delta <= 0 || delta >= 1 {
		return nil, errors.New("value of delta should be in range of (0, 1)")
	}

	w := uint(math.Ceil(2 / epsilon))
	d := uint(math.Ceil(math.Log(1-delta) / math.Log(0.5)))
	// fmt.Printf("ε: %f, δ: %f -> d: %d, w: %d\n", epsilon, delta, d, w)
	return New(d, w, topK)
}

func (s *CountMinSketch) Record(item Serializable) {
	var min uint64
	h := s.baseHashes(item.Bytes())
	for i := uint(0); i < s.d; i++ {
		loc := s.locations(h, i)
		s.m[i][loc] += 1
		if i == 0 || s.m[i][loc] < min {
			min = s.m[i][loc]
		}
	}
	s.topHeap.Push(item, min)
}

// get the two basic hash function values for data.
// Based on https://github.com/willf/bloom/blob/master/bloom.go
func (s *CountMinSketch) baseHashes(key []byte) [2]uint32 {
	s.h.Reset()
	s.h.Write(key)
	sum := s.h.Sum(nil)
	return [2]uint32{binary.BigEndian.Uint32(sum[0:4]), binary.BigEndian.Uint32(sum[4:8])}
}

func (s *CountMinSketch) locations(hashes [2]uint32, r uint) uint32 {
	return (hashes[0] + hashes[1]*uint32(r)) % uint32(s.w)
}

// Depth returns the number of hashing functions
func (s *CountMinSketch) Depth() uint {
	return s.d
}

// Width returns the size of hashing functions
func (s *CountMinSketch) Width() uint {
	return s.w
}

func (s *CountMinSketch) TopK() []*counter {
	return s.topHeap.SortedArray()
}
