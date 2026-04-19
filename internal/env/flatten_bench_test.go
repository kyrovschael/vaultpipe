package env

import (
	"fmt"
	"testing"
)

func BenchmarkFlattenMap_Wide(b *testing.B) {
	src := make(map[string]any, 100)
	for i := 0; i < 100; i++ {
		src[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	opts := DefaultFlattenOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FlattenMap(src, opts)
	}
}

func BenchmarkFlattenMap_Deep(b *testing.B) {
	src := map[string]any{
		"l1": map[string]any{
			"l2": map[string]any{
				"l3": map[string]any{
					"l4": "value",
				},
			},
		},
	}
	opts := DefaultFlattenOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FlattenMap(src, opts)
	}
}
