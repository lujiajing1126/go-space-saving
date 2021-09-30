package space_saving_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	spaceSaving "github.com/lujiajing1126/go-space-saving"
)

func prepareStream(sample string) []string {
	return strings.Split(sample, " ")
}

func Test_Summary_Stream(t *testing.T) {
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
