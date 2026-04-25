package env

import (
	"fmt"
	"testing"
)

func BenchmarkIndexSnapshot_1000Keys(b *testing.B) {
	pairs := make([]string, 1000)
	for i := range pairs {
		pairs[i] = fmt.Sprintf("KEY_%d=value_%d", i, i%50) // 50 unique values
	}
	snap := FromSlice(pairs)
	opts := DefaultIndexOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IndexSnapshot(snap, opts)
	}
}

func BenchmarkIndexSnapshot_CaseFold_500Keys(b *testing.B) {
	pairs := make([]string, 500)
	for i := range pairs {
		pairs[i] = fmt.Sprintf("KEY_%d=Value_%d", i, i%25)
	}
	snap := FromSlice(pairs)
	opts := DefaultIndexOptions()
	opts.CaseFold = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IndexSnapshot(snap, opts)
	}
}
