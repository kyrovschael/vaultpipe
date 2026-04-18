// Package rotate provides secret rotation detection and reload triggering.
package rotate

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// SecretFetcher retrieves a map of secret key→value pairs.
type SecretFetcher func(ctx context.Context) (map[string]string, error)

// ChangeHandler is called when a change in secrets is detected.
type ChangeHandler func(updated map[string]string)

// Watcher polls secrets and invokes a handler when they change.
type Watcher struct {
	fetch    SecretFetcher
	onChange ChangeHandler
	interval time.Duration

	mu      sync.Mutex
	lastSig string
}

// NewWatcher creates a Watcher that polls every interval.
func NewWatcher(fetch SecretFetcher, onChange ChangeHandler, interval time.Duration) *Watcher {
	return &Watcher{
		fetch:    fetch,
		onChange: onChange,
		interval: interval,
	}
}

// Run starts the polling loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			secrets, err := w.fetch(ctx)
			if err != nil {
				continue
			}
			sig := signature(secrets)
			w.mu.Lock()
			changed := sig != w.lastSig
			w.lastSig = sig
			w.mu.Unlock()
			if changed {
				w.onChange(secrets)
			}
		}
	}
}

// signature returns a deterministic hash of the secret map.
func signature(secrets map[string]string) string {
	h := sha256.New()
	for k, v := range secrets {
		fmt.Fprintf(h, "%s=%s;", k, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
