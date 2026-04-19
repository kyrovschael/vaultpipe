package env

import (
	"fmt"
	"testing"
)

func BenchmarkRenameSnapshot_500Keys(b *testing.B) {
	s := make(Snapshot, 500)
	m := make(map[string]string, 250)
	for i := 0; i < 500; i++ {
		k := fmt.Sprintf("KEY_%d", i)
		s[k] = fmt.Sprintf("val_%d", i)
		if i < 250 {
			m[k] = fmt.Sprintf("RENAMED_%d", i)
		}
	}
	opts := DefaultRenameOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenameSnapshot(s, m, opts)
	}
}

func BenchmarkRenameSnapshot_DropUnmapped_500Keys(b *testing.B) {
	s := make(Snapshot, 500)
	m := make(map[string]string, 50)
	for i := 0; i < 500; i++ {
		k := fmt.Sprintf("KEY_%d", i)
		s[k] = fmt.Sprintf("val_%d", i)
		if i < 50 {
			m[k] = fmt.Sprintf("RENAMED_%d", i)
		}
	}
	opts := RenameOptions{DropUnmapped: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenameSnapshot(s, m, opts)
	}
}
