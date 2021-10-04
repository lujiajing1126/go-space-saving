package benchmark

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

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

func Benchmark_SpaceSaving_Insert(b *testing.B) {
	ss, err := spaceSaving.NewStreamSummary(0.01)
	if err != nil {
		b.Error("fail to create stream summary", err)
	}
	f, err := os.Open("../pkg/testdata/lorem_ipsum.txt")
	if err != nil {
		b.Error("fail to open file", err)
	}
	data, _ := io.ReadAll(f)
	stream := prepareStream(string(data))
	b.ResetTimer()
	for i, l := 0, len(stream); i < b.N; i++ {
		ss.Record(stream[i%l])
	}
}
