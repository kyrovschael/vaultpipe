package env

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatchSnapshot_DetectsChange(t *testing.T) {
	var callCount atomic.Int32
	var lastSnap Snapshot

	gen := func() int { return int(callCount.Load()) }
	source := func() (Snapshot, error) {
		return Snapshot{"KEY": "value" + string(rune('0'+gen()))}, nil
	}

	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(s Snapshot) {
			callCount.Add(1)
			lastSnap = s
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	WatchSnapshot(ctx, source, opts) //nolint:errcheck

	if callCount.Load() == 0 {
		t.Fatal("expected at least one change callback")
	}
	if lastSnap == nil {
		t.Fatal("lastSnap should not be nil")
	}
}

func TestWatchSnapshot_NoChangeNoCallback(t *testing.T) {
	var callCount atomic.Int32

	source := func() (Snapshot, error) {
		return Snapshot{"KEY": "stable"}, nil
	}

	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(s Snapshot) { callCount.Add(1) },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	WatchSnapshot(ctx, source, opts) //nolint:errcheck

	if callCount.Load() != 0 {
		t.Fatalf("expected no callbacks, got %d", callCount.Load())
	}
}

func TestWatchSnapshot_SourceErrorSkipped(t *testing.T) {
	var callCount atomic.Int32
	fail := true

	source := func() (Snapshot, error) {
		if fail {
			fail = false
			return nil, errors.New("transient")
		}
		return Snapshot{"K": "v"}, nil
	}

	opts := WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(s Snapshot) { callCount.Add(1) },
	}

	// First call to source returns error — WatchSnapshot should propagate it.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := WatchSnapshot(ctx, source, opts)
	if err == nil {
		t.Fatal("expected error from initial source call")
	}
}

func TestSnapshotsEqual(t *testing.T) {
	a := Snapshot{"A": "1", "B": "2"}
	b := Snapshot{"A": "1", "B": "2"}
	if !snapshotsEqual(a, b) {
		t.Fatal("expected equal")
	}
	b["A"] = "changed"
	if snapshotsEqual(a, b) {
		t.Fatal("expected not equal")
	}
}
