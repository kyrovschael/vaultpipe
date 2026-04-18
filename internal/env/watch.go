package env

import (
	"context"
	"time"
)

// WatchOptions configures the environment snapshot watcher.
type WatchOptions struct {
	// Interval between polls.
	Interval time.Duration
	// OnChange is called with the new snapshot whenever a change is detected.
	OnChange func(snap Snapshot)
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions(onChange func(Snapshot)) WatchOptions {
	return WatchOptions{
		Interval: 5 * time.Second,
		OnChange: onChange,
	}
}

// WatchSnapshot polls a source function at the configured interval and calls
// OnChange whenever the returned snapshot differs from the previous one.
// It blocks until ctx is cancelled.
func WatchSnapshot(ctx context.Context, source func() (Snapshot, error), opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = 5 * time.Second
	}

	current, err := source()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			next, err := source()
			if err != nil {
				continue
			}
			if !snapshotsEqual(current, next) {
				current = next
				if opts.OnChange != nil {
					opts.OnChange(current)
				}
			}
		}
	}
}

func snapshotsEqual(a, b Snapshot) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
