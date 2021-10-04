package spacesaving_test

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	spaceSaving "github.com/lujiajing1126/go-space-saving/pkg/spacesaving"
)

var (
	WhiteSpaces  = regexp.MustCompile("\\s+")
	Punctuations = regexp.MustCompile("[\\.,]")
)

func prepareStream(sample string) []string {
	sample = Punctuations.ReplaceAllString(sample, "")
	sample = WhiteSpaces.ReplaceAllString(sample, " ")
	return strings.Split(sample, " ")
}

func Test_SummaryStream(t *testing.T) {
	tester := assert.New(t)
	ss, err := spaceSaving.NewStreamSummary(0.01)
	tester.NoError(err)
	tester.NotNil(ss)

	// Data Stream
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	stream := prepareStream(sample)
	for _, singleItem := range stream {
		ss.Record(singleItem)
	}
	top3Counters := ss.TopK(3)
	tester.Len(top3Counters, 3)
	var top3Strings []string
	for _, top3item := range top3Counters {
		top3Strings = append(top3Strings, top3item.Val().(string))
	}
	tester.Equal(top3Strings, []string{"g", "f", "e"})
}

func Test_SummaryStream_givenItemsSeenInTheStream_Is_LargerThanCapacity(t *testing.T) {
	tester := assert.New(t)
	// Init with 100 counters
	ss, err := spaceSaving.NewStreamSummary(0.01)
	tester.NoError(err)
	tester.NotNil(ss)

	// Data Stream
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	stream := prepareStream(sample)
	for _, singleItem := range stream {
		ss.Record(singleItem)
	}

	// Record skewed data stream
	for i := 0; i < 200; i++ {
		ss.Record("a")
	}

	top3Counters := ss.TopK(3)
	tester.Len(top3Counters, 3)
	var top3Strings []string
	for _, top3item := range top3Counters {
		top3Strings = append(top3Strings, top3item.Val().(string))
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

	// 500 counters with error range of ~9 words (4623 * 0.002)
	ss, err := spaceSaving.NewStreamSummary(0.002)
	tester.NoError(err)
	tester.NotNil(ss)

	stream := prepareStream(string(sample))
	for _, singleItem := range stream {
		ss.Record(singleItem)
	}

	topCounters := ss.TopK(5)
	tester.Len(topCounters, 5)
	var topStrings []string
	for _, topitem := range topCounters {
		topStrings = append(topStrings, topitem.Val().(string))
	}
	tester.Contains(topStrings, "et")
}
