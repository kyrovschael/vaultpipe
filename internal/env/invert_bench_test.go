package env

import (
	"fmt"
	"testing"
)

func BenchmarkInvertSnapshot_500UniqueValues(b *testing.B) {
	src := make(Snapshot, 500)
	for i := 0; i < 500; i++ {
		src[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("val_%d", i)
	}
	opts := DefaultInvertOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InvertSnapshot(src, opts)
	}
}

func BenchmarkInvertSnapshot_500DuplicateValues(b *testing.B) {
	// All keys share the same value — worst-case duplicate handling.
	src := make(Snapshot, 500)
	for i := 0; i < 500; i++ {
		src[fmt.Sprintf("KEY_%d", i)] = "shared_value"
	}
	opts := DefaultInvertOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InvertSnapshot(src, opts)
	}
}
