package limit

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_InvalidConcurrent(t *testing.T) {
	_, err := New(Config{MaxConcurrent: 0, RequestTimeout: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxConcurrent=0")
	}
}

func TestNew_InvalidTimeout(t *testing.T) {
	_, err := New(Config{MaxConcurrent: 1, RequestTimeout: 0})
	if err == nil {
		t.Fatal("expected error for RequestTimeout=0")
	}
}

func TestAcquire_Release(t *testing.T) {
	l, err := New(Config{MaxConcurrent: 2, RequestTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	release, err := l.Acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	release()
}

func TestAcquire_CancelledContext(t *testing.T) {
	l, _ := New(Config{MaxConcurrent: 1, RequestTimeout: time.Second})
	// Fill the slot.
	release, _ := l.Acquire(context.Background())
	defer release()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := l.Acquire(ctx)
	if err == nil {
		t.Fatal("expected context error when slot unavailable")
	}
}

func TestAcquire_MaxConcurrency(t *testing.T) {
	const max = 4
	l, _ := New(Config{MaxConcurrent: max, RequestTimeout: time.Second})

	var active int64
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release, err := l.Acquire(context.Background())
			if err != nil {
				return
			}
			defer release()
			v := atomic.AddInt64(&active, 1)
			if v > max {
				t.Errorf("active=%d exceeds max=%d", v, max)
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&active, -1)
		}()
	}
	wg.Wait()
}

func TestWithTimeout(t *testing.T) {
	l, _ := New(Config{MaxConcurrent: 1, RequestTimeout: 100 * time.Millisecond})
	ctx, cancel := l.WithTimeout(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) > 100*time.Millisecond {
		t.Error("deadline further than expected")
	}
}
