package env

import (
	"fmt"
	"testing"
)

func BenchmarkWindowSnapshot_1000Keys(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := range s {
		s[i] = Entry{Key: fmt.Sprintf("KEY_%04d", i), Value: fmt.Sprintf("val%d", i)}
	}
	opts := WindowOptions{From: "KEY_0200", To: "KEY_0800"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WindowSnapshot(s, opts)
	}
}

func BenchmarkWindowSnapshot_CaseFold_1000Keys(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := range s {
		s[i] = Entry{Key: fmt.Sprintf("key_%04d", i), Value: fmt.Sprintf("val%d", i)}
	}
	opts := WindowOptions{From: "KEY_0200", To: "KEY_0800", CaseFold: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WindowSnapshot(s, opts)
	}
}
