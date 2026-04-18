package rotate_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/rotate"
)

func TestWatcher_DetectsChange(t *testing.T) {
	calls := int32(0)
	generation := int32(0)

	fetch := func(ctx context.Context) (map[string]string, error) {
		g := atomic.LoadInt32(&generation)
		if g == 0 {
			return map[string]string{"KEY": "v1"}, nil
		}
		return map[string]string{"KEY": "v2"}, nil
	}

	handler := func(_ map[string]string) {
		atomic.AddInt32(&calls, 1)
	}

	w := rotate.NewWatcher(fetch, handler, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go w.Run(ctx) //nolint:errcheck

	time.Sleep(50 * time.Millisecond)
	atomic.StoreInt32(&generation, 1)
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&calls) < 1 {
		t.Fatal("expected handler to be called after secret change")
	}
}

func TestWatcher_NoChangeNoCall(t *testing.T) {
	calls := int32(0)

	fetch := func(ctx context.Context) (map[string]string, error) {
		return map[string]string{"KEY": "stable"}, nil
	}
	handler := func(_ map[string]string) {
		atomic.AddInt32(&calls, 1)
	}

	w := rotate.NewWatcher(fetch, handler, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	go w.Run(ctx) //nolint:errcheck
	time.Sleep(130 * time.Millisecond)

	// First poll sets baseline; subsequent identical polls must not fire.
	if n := atomic.LoadInt32(&calls); n > 1 {
		t.Fatalf("expected at most 1 call (initial), got %d", n)
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	fetch := func(ctx context.Context) (map[string]string, error) {
		return map[string]string{"K": "v"}, nil
	}
	w := rotate.NewWatcher(fetch, func(_ map[string]string) {}, 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected non-nil error on cancelled context")
	}
}
