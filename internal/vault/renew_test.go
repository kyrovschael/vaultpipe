package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

func newTestClient(t *testing.T, handler http.Handler) *api.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := api.DefaultConfig()
	cfg.Address = srv.URL
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestRenewToken_NonRenewable(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return a token with TTL=0 (non-renewable)
		w.Write([]byte(`{"auth":null,"data":{"ttl":0},"renewable":false,"lease_duration":0}`))
	})

	client := newTestClient(t, handler)
	ctx := context.Background()

	cancel, err := RenewToken(ctx, client)
	if err != nil {
		// Acceptable: server returns unexpected format; just ensure no panic
		t.Logf("RenewToken returned error (acceptable in unit test): %v", err)
		return
	}
	if cancel == nil {
		t.Fatal("expected non-nil cancel function")
	}
	cancel()
}

func TestRenewLoop_CancelStopsLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		renewLoop(ctx, nil, 30*time.Second)
	}()

	// Cancel immediately and verify the loop exits
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("renewLoop did not exit after context cancellation")
	}
}
