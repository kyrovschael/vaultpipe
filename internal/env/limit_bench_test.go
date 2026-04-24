package env

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkLimitSnapshot_1000Keys_MaxKeys(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := range s {
		s[i] = Entry{Key: fmt.Sprintf("KEY_%04d", i), Value: fmt.Sprintf("value_%d", i)}
	}
	opts := LimitOptions{MaxKeys: 500}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LimitSnapshot(s, opts)
	}
}

func BenchmarkLimitSnapshot_1000Keys_MaxValueLen(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := range s {
		v := fmt.Sprintf("value_%d", i)
		if i%3 == 0 {
			v = strings.Repeat("x", 512)
		}
		s[i] = Entry{Key: fmt.Sprintf("KEY_%04d", i), Value: v}
	}
	opts := LimitOptions{MaxValueLen: 64}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LimitSnapshot(s, opts)
	}
}
