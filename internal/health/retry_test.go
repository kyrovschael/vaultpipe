package health_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/health"
)

func TestWaitUntilHealthy_SucceedsOnSecondAttempt(t *testing.T) {
	var calls atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if calls.Add(1) < 2 {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"initialized":true,"sealed":true,"standby":false,"version":"1.15.0"}`))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false,"standby":false,"version":"1.15.0"}`))
	})
	checker := newTestChecker(t, handler)
	cfg := health.RetryConfig{MaxAttempts: 5, Interval: 10 * time.Millisecond}
	s := health.WaitUntilHealthy(context.Background(), checker, cfg)
	if !s.Healthy {
		t.Fatalf("expected healthy after retry, err: %v", s.Error)
	}
	if calls.Load() < 2 {
		t.Errorf("expected at least 2 calls, got %d", calls.Load())
	}
}

func TestWaitUntilHealthy_ExhaustsAttempts(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":true,"standby":false,"version":"1.15.0"}`))
	})
	checker := newTestChecker(t, handler)
	cfg := health.RetryConfig{MaxAttempts: 3, Interval: 5 * time.Millisecond}
	s := health.WaitUntilHealthy(context.Background(), checker, cfg)
	if s.Healthy {
		t.Fatal("expected unhealthy after exhausting attempts")
	}
}

func TestWaitUntilHealthy_ContextCancelled(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":true,"standby":false,"version":"1.15.0"}`))
	})
	checker := newTestChecker(t, handler)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cfg := health.RetryConfig{MaxAttempts: 10, Interval: time.Second}
	s := health.WaitUntilHealthy(ctx, checker, cfg)
	if s.Healthy {
		t.Fatal("expected unhealthy when context cancelled")
	}
}
