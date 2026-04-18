// Package limit provides rate-limiting and concurrency controls for
// Vault secret fetches to avoid overwhelming the Vault server.
package limit

import (
	"context"
	"fmt"
	"time"
)

// Config holds rate-limiting parameters.
type Config struct {
	// MaxConcurrent is the maximum number of simultaneous Vault requests.
	MaxConcurrent int
	// RequestTimeout is the per-request deadline.
	RequestTimeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxConcurrent:  8,
		RequestTimeout: 10 * time.Second,
	}
}

// Limiter gates concurrent access to Vault.
type Limiter struct {
	sem     chan struct{}
	timeout time.Duration
}

// New creates a Limiter from cfg.
func New(cfg Config) (*Limiter, error) {
	if cfg.MaxConcurrent <= 0 {
		return nil, fmt.Errorf("limit: MaxConcurrent must be > 0, got %d", cfg.MaxConcurrent)
	}
	if cfg.RequestTimeout <= 0 {
		return nil, fmt.Errorf("limit: RequestTimeout must be > 0, got %s", cfg.RequestTimeout)
	}
	return &Limiter{
		sem:     make(chan struct{}, cfg.MaxConcurrent),
		timeout: cfg.RequestTimeout,
	}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns a cancel func that must be called to release the slot.
func (l *Limiter) Acquire(ctx context.Context) (release func(), err error) {
	select {
	case l.sem <- struct{}{}:
		return func() { <-l.sem }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// WithTimeout wraps ctx with the configured per-request timeout.
func (l *Limiter) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, l.timeout)
}
