package env

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkChainSource_10Sources(b *testing.B) {
	sources := make([]SourceFunc, 10)
	for i := range sources {
		m := make(map[string]string, 100)
		for j := 0; j < 100; j++ {
			m[fmt.Sprintf("KEY_%d_%d", i, j)] = fmt.Sprintf("val_%d", j)
		}
		sources[i] = MapSource(m)
	}
	chain := ChainSource(sources...)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chain(ctx)
	}
}

func BenchmarkStaticSource_Clone(b *testing.B) {
	m := make(Snapshot, 500)
	for i := 0; i < 500; i++ {
		m[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	src := StaticSource(m)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = src(ctx)
	}
}
