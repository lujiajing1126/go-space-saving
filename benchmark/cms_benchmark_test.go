package benchmark

import (
	"io"
	"os"
	"testing"

	"github.com/lujiajing1126/go-space-saving/pkg/cms"
)

type StringRepl string

func (sr StringRepl) Bytes() []byte {
	return []byte(sr)
}

func Benchmark_CountMinSketch_Insert_Top5(b *testing.B) {
	cm, err := cms.NewWithEstimates(0.005, 0.999, 10)
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
		cm.Record(StringRepl(stream[i%l]))
	}
}
