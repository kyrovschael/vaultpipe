package env

import "testing"

func BenchmarkSubsetSnapshot_100of1000(b *testing.B) {
	src := make(Snapshot, 1000)
	keys := make([]string, 100)
	for i := 0; i < 1000; i++ {
		k := "KEY_" + itoa(i)
		src[k] = "value"
		if i < 100 {
			keys[i] = k
		}
	}
	opts := SubsetOptions{Keys: keys, IgnoreMissing: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SubsetSnapshot(src, opts)
	}
}
