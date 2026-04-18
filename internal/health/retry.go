package health

import (
	"context"
	"time"
)

// RetryConfig controls retry behaviour for WaitUntilHealthy.
type RetryConfig struct {
	MaxAttempts int
	Interval    time.Duration
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 10,
		Interval:    2 * time.Second,
	}
}

// WaitUntilHealthy polls Vault until it is healthy or the context is cancelled.
// It returns the last Status observed.
func WaitUntilHealthy(ctx context.Context, checker *Checker, cfg RetryConfig) Status {
	var last Status
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		last = checker.Check(ctx)
		if last.Healthy {
			return last
		}
		select {
		case <-ctx.Done():
			last.Error = ctx.Err()
			return last
		case <-time.After(cfg.Interval):
		}
	}
	return last
}
