package env

import (
	"fmt"
	"testing"
)

func BenchmarkTransformSnapshot_1000Keys_UpperKeys(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := 0; i < 1000; i++ {
		s[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	opts := UpperKeys()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TransformSnapshot(s, opts)
	}
}

func BenchmarkTransformSnapshot_1000Keys_TrimValues(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := 0; i < 1000; i++ {
		s[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("  value_%d  ", i)
	}
	opts := TrimValues()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TransformSnapshot(s, opts)
	}
}
