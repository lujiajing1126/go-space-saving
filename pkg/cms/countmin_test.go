package cms

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StringRepl string

func (sr StringRepl) Bytes() []byte {
	return []byte(sr)
}

var (
	WhiteSpaces  = regexp.MustCompile("\\s+")
	Punctuations = regexp.MustCompile("[\\.,]")
)

func prepareStream(sample string) []string {
	sample = Punctuations.ReplaceAllString(sample, "")
	sample = WhiteSpaces.ReplaceAllString(sample, " ")
	return strings.Split(sample, " ")
}

func Test_CountMinSketch(t *testing.T) {
	tester := assert.New(t)
	ss, err := New(10, 4, 3)
	tester.NoError(err)
	tester.NotNil(ss)

	// Data Stream
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	stream := prepareStream(sample)
	for _, singleItem := range stream {
		ss.Record(StringRepl(singleItem))
	}
	top3Counters := ss.TopK()
	tester.Len(top3Counters, 3)
	var top3Strings []string
	for _, top3item := range top3Counters {
		top3Strings = append(top3Strings, string(top3item.Val().(StringRepl)))
	}
	tester.Equal(top3Strings, []string{"g", "f", "e"})
}

func Test_CountMinSketch_givenItemsSeenInTheStream_Is_LargerThanCapacity(t *testing.T) {
	tester := assert.New(t)
	// Init with 100 counters
	ss, err := New(10, 4, 3)
	tester.NoError(err)
	tester.NotNil(ss)

	// Data Stream
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	stream := prepareStream(sample)
	for _, singleItem := range stream {
		ss.Record(StringRepl(singleItem))
	}

	// Record skewed data stream
	for i := 0; i < 200; i++ {
		ss.Record(StringRepl("a"))
	}

	top3Counters := ss.TopK()
	tester.Len(top3Counters, 3)
	var top3Strings []string
	for _, top3item := range top3Counters {
		top3Strings = append(top3Strings, string(top3item.Val().(StringRepl)))
	}
	tester.Equal(top3Strings, []string{"a", "g", "f"})
}

// Generated 50 paragraphs, 4623 words, 31168 bytes of Lorem Ipsum
// With 1050 uniq words
// $ tr -c '[:alnum:]' '[\n*]' < testdata/lorem_ipsum.txt | sort | uniq -c | sort -nr | head  -
func Test_Lorem_Ipsum(t *testing.T) {
	tester := assert.New(t)

	f, err := os.Open("../testdata/lorem_ipsum.txt")
	tester.NoError(err)
	tester.NotNil(f)
	sample, err := ioutil.ReadAll(f)
	tester.NoError(err)

	// 500 counters with error range of ~9 words (4623 * 0.001)
	ss, err := NewWithEstimates(0.01, 0.99, 5)
	tester.NoError(err)
	tester.NotNil(ss)

	stream := prepareStream(string(sample))
	for _, singleItem := range stream {
		ss.Record(StringRepl(singleItem))
	}

	topCounters := ss.TopK()
	tester.Len(topCounters, 5)
	var topStrings []string
	for _, topitem := range topCounters {
		topStrings = append(topStrings, string(topitem.Val().(StringRepl)))
	}
	tester.Contains(topStrings, "et")
}
