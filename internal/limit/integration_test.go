package limit_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/limit"
)

// TestDefaultConfig_Usable ensures the default config produces a working Limiter.
func TestDefaultConfig_Usable(t *testing.T) {
	cfg := limit.DefaultConfig()
	l, err := limit.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const workers = 16
	var completed int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := l.WithTimeout(context.Background())
			defer cancel()

			release, err := l.Acquire(ctx)
			if err != nil {
				t.Errorf("acquire failed: %v", err)
				return
			}
			defer release()
			// Simulate a short Vault fetch.
			time.Sleep(2 * time.Millisecond)
			atomic.AddInt64(&completed, 1)
		}()
	}

	wg.Wait()
	if got := atomic.LoadInt64(&completed); got != workers {
		t.Errorf("completed=%d want %d", got, workers)
	}
}
