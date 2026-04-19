package env

import (
	"fmt"
	"testing"
)

func BenchmarkPipeline_1000Keys_ThreeSteps(b *testing.B) {
	pairs := make([]string, 1000)
	for i := range pairs {
		pairs[i] = fmt.Sprintf("key_%d=value_%d", i, i)
	}
	src := FromSlice(pairs)

	p := NewPipeline(
		Lift(UpperKeys),
		Lift(TrimValues),
		Lift(func(s Snapshot) Snapshot { return s }),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Run(src)
	}
}
