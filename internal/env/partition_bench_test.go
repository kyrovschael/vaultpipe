package env

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkPartitionSnapshot_1000Keys(b *testing.B) {
	src := make(Snapshot, 1000)
	for i := 0; i < 500; i++ {
		src[fmt.Sprintf("DB_%04d", i)] = fmt.Sprintf("value%d", i)
	}
	for i := 500; i < 1000; i++ {
		src[fmt.Sprintf("APP_%04d", i)] = fmt.Sprintf("value%d", i)
	}
	opts := DefaultPartitionOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PartitionSnapshot(src, opts, func(k, _ string) bool {
			return strings.HasPrefix(k, "DB_")
		})
	}
}

func BenchmarkPartitionSnapshot_CaseFold_1000Keys(b *testing.B) {
	src := make(Snapshot, 1000)
	for i := 0; i < 1000; i++ {
		src[fmt.Sprintf("key_%04d", i)] = fmt.Sprintf("val%d", i)
	}
	opts := DefaultPartitionOptions()
	opts.CaseFold = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PartitionSnapshot(src, opts, func(k, _ string) bool {
			return strings.HasPrefix(k, "KEY_0")
		})
	}
}
