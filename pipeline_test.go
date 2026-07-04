package main

import (
	"bytes"
	"testing"
)

func BenchmarkConcurrent(b *testing.B) {
	if len(csvBytes) == 0 {
		b.Skip("commits.csv not found")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(csvBytes)
		_, _ = RunPipeline(reader, 0)
	}
}
